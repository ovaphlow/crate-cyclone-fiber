package genericimplementation

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func EndpointGet(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	result, err := retrieve(&schema, &table)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(result) > 0 {
		for i, item := range result {
			for k, v := range item {
				switch val := v.(type) {
				case int64:
					item["_"+k] = strconv.FormatInt(val, 10)
				}
			}
			result[i] = item
		}
		return c.JSON(result)
	}
	return c.SendString("[]")
}
