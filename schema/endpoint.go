package schema

import (
	"ovaphlow/cratecyclone/utilities"

	"github.com/gofiber/fiber/v2"
)

func EndpointGet(c *fiber.Ctx) error {
	schemas, err := retrieveSchemas()
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(schemas) > 0 {
		return c.JSON(schemas)
	}
	return c.JSON([]string{})
}
