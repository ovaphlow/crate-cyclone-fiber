package subscriber

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"ovaphlow/cratecyclone/configuration"
	"ovaphlow/cratecyclone/utility"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func GetWithParams(c *fiber.Ctx) error {
	c.Set(configuration.HeaderAPIVersion, "2024-01-06")
	id := c.Params("id", "")
	uuid := c.Params("uuid", "")
	if id == "" || uuid == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	id_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberById(id_, uuid)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber == nil {
		return c.Status(404).JSON(fiber.Map{"message": "用户不存在"})
	}
	return c.JSON(fiber.Map{
		"id":     subscriber.Id,
		"email":  subscriber.Email,
		"name":   subscriber.Name,
		"phone":  subscriber.Phone,
		"tags":   subscriber.Tags,
		"detail": subscriber.Detail,
		"time":   subscriber.Time,
		"_id":    strconv.FormatInt(subscriber.Id, 10),
	})
}

func RefreshJwt(c *fiber.Ctx) error {
	c.Set(configuration.HeaderAPIVersion, "2024-01-06")
	type Body struct {
		Token string `json:"token"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Token == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
	})
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utility.Slogger.Error("token claims is not jwt.MapClaims")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if claims["exp"] == nil {
		utility.Slogger.Error("token claims exp is nil")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return c.Status(401).JSON(fiber.Map{"message": "token 已过期"})
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")))
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

func SignIn(c *fiber.Ctx) error {
	c.Set(configuration.HeaderAPIVersion, "2024-01-06")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	jwtKey := []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", ""))
	type SignInBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body SignInBody
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Username == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberByUsername(body.Username)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber == nil {
		utility.Slogger.Error("用户不存在")
		return c.Status(401).JSON(fiber.Map{"message": "用户名或密码错误"})
	}
	var detail map[string]interface{}
	if err := json.Unmarshal([]byte(subscriber.Detail), &detail); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	salt, ok := detail["salt"].(string)
	if !ok {
		utility.Slogger.Error("salt is not string")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	key := []byte(salt)
	r := hmac.New(sha256.New, key)
	r.Write([]byte(body.Password))
	sha := hex.EncodeToString(r.Sum(nil))
	if sha != detail["sha"] {
		return c.Status(401).JSON(fiber.Map{"message": "用户名或密码错误"})
	}
	var state map[string]interface{}
	if err := json.Unmarshal([]byte(subscriber.State), &state); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    "crate",
		Subject:   strconv.FormatInt(subscriber.Id, 10),
		Id:        state["uuid"].(string),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

func SignUp(c *fiber.Ctx) error {
	c.Set(configuration.HeaderAPIVersion, "2024-01-06")
	type SignUpBody struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	var body SignUpBody
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if (body.Email == "" && body.Name == "" && body.Phone == "") || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberByUsername(body.Email)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber != nil {
		return c.Status(401).JSON(fiber.Map{"message": "用户名已存在"})
	}
	subscriber = &Subscriber{
		Email: body.Email,
		Name:  body.Name,
		Phone: body.Phone,
		Tags:  "[]",
	}
	bytes := make([]byte, 8)
	_, err = rand.Read(bytes)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	key := []byte(hex.EncodeToString(bytes))
	r := hmac.New(sha256.New, key)
	r.Write([]byte(body.Password))
	sha := hex.EncodeToString(r.Sum(nil))
	subscriber.Detail = fmt.Sprintf(`{"salt": "%s", "sha": "%s", "uuid": "%s"}`, key, sha, uuid.New())
	if err := repoCreateSubscriber(subscriber); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"message": "注册成功"})
}

func ValidateToken(c *fiber.Ctx) error {
	type Body struct {
		Token string `json:"token"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Token == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
	})
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utility.Slogger.Error("token claims is not jwt.MapClaims")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if claims["exp"] == nil {
		utility.Slogger.Error("token claims exp is nil")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return c.Status(401).JSON(fiber.Map{"message": "token 已过期"})
	}
	userId, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	result, err := repoRetrieveSubscriberById(userId, claims["jti"].(string))
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if result == nil {
		return c.Status(401).JSON(fiber.Map{"message": "用户不存在"})
	}
	claims["email"] = result.Email
	claims["name"] = result.Name
	claims["phone"] = result.Phone
	claims["tags"] = result.Tags
	return c.JSON(claims)
}
