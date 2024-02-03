package schema

import (
	"ovaphlow/cratecyclone/utilities"

	"github.com/gofiber/fiber/v2"
)

func EndpointDeleteSchema(c *fiber.Ctx) error {
	name := c.Params("name")
	err := removeSchema(name)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(200)
}

func EndpointGetSchema(c *fiber.Ctx) error {
	schemas, err := retrieveSchemas()
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(fiber.Map{"schemas": schemas})
}

func EndpointPostSchema(c *fiber.Ctx) error {
	type body struct {
		Schema string `json:"schema"`
	}
	var b body
	err := c.BodyParser(&b)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	err = createSchema(b.Schema)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(201)
}

func EndpointPutSchema(c *fiber.Ctx) error {
	type body struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	var b body
	err := c.BodyParser(&b)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	err = updateSchema(b.Current, b.New)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(200)
}
