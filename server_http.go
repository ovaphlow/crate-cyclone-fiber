package main

import (
	"log"
	"os"
	"ovaphlow/cratecyclone/configuration"
	"ovaphlow/cratecyclone/schema"
	"ovaphlow/cratecyclone/subscriber"
	"ovaphlow/cratecyclone/utility"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt"
)

func HTTPServe(addr string) {
	utility.InitPostgres()

	app := fiber.New(fiber.Config{
		// Prefork 模式与 gRPC 服务冲突
		// Prefork:   true,
		BodyLimit: 16 * 1024 * 1024,
	})

	app.Use(compress.New())

	app.Use(cors.New())

	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(func(c *fiber.Ctx) error {
		utility.Slogger.Info(c.Path(), "method", c.Method(), "query", c.Queries(), "ip", c.IP())
		return c.Next()
	})

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set(configuration.HeaderAPIVersion, "2024-02-03")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		for _, item := range configuration.PublicUris {
			match, _ := regexp.MatchString(item, c.Path())
			if match {
				return c.Next()
			}
		}
		auth := c.Get("Authorization")
		auth = strings.Replace(auth, "Bearer ", "", 1)
		token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
		})
		if err != nil {
			utility.Slogger.Error(err.Error())
			return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		}
		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		}
		return c.Next()
	})

	app.Get("/cyclone-api/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hola el mondo",
		})
	})

	app.Post("/cyclone-api/subscriber/refresh-jwt", subscriber.RefreshJwt)
	app.Post("/cyclone-api/subscriber/sign-in", subscriber.SignIn)
	app.Post("/cyclone-api/subscriber/sign-up", subscriber.SignUp)
	app.Post("/cyclone-api/subscriber/validate-token", subscriber.ValidateToken)
	app.Get("/cyclone-api/subscriber/:uuid/:id", subscriber.GetWithParams)

	app.Get("/cyclone-api/db-schema", schema.GetSchema)
	app.Get("/cyclone-api/:schema/db-table", schema.GetTable)
	app.Get("/cyclone-api/:schema/:table", schema.Get)
	app.Get("/cyclone-api/:schema/:table/:uuid/:id", schema.GetWithParams)
	app.Post("/cyclone-api/:schema/:table", schema.Post)
	app.Put("/cyclone-api/:schema/:table/:uuid/:id", schema.Put)
	app.Delete("/cyclone-api/:schema/:table/:uuid/:id", schema.Delete)

	app.Get("/cyclone-api/get/db-schema", schema.GetSchema)
	app.Get("/cyclone-api/get/db-table/:schema", schema.GetTable)

	log.Fatal(app.Listen(addr))
}
