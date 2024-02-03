package genericimplementation

import "github.com/gofiber/fiber/v2"

func EndpointGet(c *fiber.Ctx) error {
	result, err := retrieve()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(result) > 0 {
		return c.JSON(result)
	}
	return c.SendString("[]")
}
