package table

import "github.com/gofiber/fiber/v2"

func EndpointGet(c *fiber.Ctx) error {
	schema := c.Query("schema", "")
	if schema == "" {
		return c.Status(400).JSON(fiber.Map{"message": "错误的 schema 值"})
	}
	tables, err := retrieveTables(&schema)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(tables) > 0 {
		return c.JSON(tables)
	}
	return c.JSON([]string{})
}
