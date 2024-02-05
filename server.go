package main

import (
	"log"
	"os"
	"ovaphlow/cratecyclone/configurations"
	gi "ovaphlow/cratecyclone/generic-implementation"
	"ovaphlow/cratecyclone/schema"
	"ovaphlow/cratecyclone/table"
	"ovaphlow/cratecyclone/utilities"
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

func Serve(addr string) {
	utilities.InitPostgres()

	app := fiber.New(fiber.Config{
		Prefork:   true,
		BodyLimit: 16 * 1024 * 1024,
	})

	app.Use(compress.New())

	app.Use(cors.New())

	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(func(c *fiber.Ctx) error {
		utilities.Slogger.Info(c.Path(), "method", c.Method(), "query", c.Queries(), "ip", c.IP())
		return c.Next()
	})

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set(configurations.HeaderAPIVersion, "2024-02-03")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		for _, item := range configurations.PublicUris {
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
			utilities.Slogger.Error(err.Error())
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

	app.Get("/cyclone-api/db-schema", schema.EndpointGet)

	app.Get("/cyclone-api/db-table", table.EndpointGet)

	app.Get("/cyclone-api/:schema/:table", gi.EndpointGet)
	app.Post("/cyclone-api/:schema/:table", gi.EndpointPost)
	app.Put("/cyclone-api/:schema/:table/:id", gi.EndpointPut)
	app.Delete("/cyclone-api/:schema/:table/:id", gi.EndpointDelete)

	app.Get("/cyclone-api/get/db-schema", schema.EndpointGet)

	app.Get("/cyclone-api/get/db-table/:schema", table.EndpointGet)

	log.Fatal(app.Listen(addr))
}
