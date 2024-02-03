package table

import "github.com/gofiber/fiber/v2"

func EndpointGet(c *fiber.Ctx) error {
	schema := c.Params("schema")
	tables, err := retrieveTables(schema)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(tables) > 0 {
		return c.JSON(tables)
	}
	return c.JSON([]string{})
}
