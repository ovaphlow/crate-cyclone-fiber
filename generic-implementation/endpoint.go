package genericimplementation

import (
	"ovaphlow/cratecyclone/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func EndpointDelete(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	id := c.Params("id")
	uuid := c.Params("uuid")
	if err := remove(&schema, &table, &id, &uuid); err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(204)
}

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

func EndpointGetWithParams(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	uuid := c.Params("uuid")
	id := c.Params("id")
	result, err := retrieveByID(&schema, &table, &id, &uuid)
	if err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if result == nil {
		return c.SendStatus(404)
	}
	for k, v := range result {
		switch val := v.(type) {
		case int64:
			result["_"+k] = strconv.FormatInt(val, 10)
		}
	}
	return c.JSON(result)
}

func EndpointPost(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	if err := create(&schema, &table, body); err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(201)
}

func EndpointPut(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	uuid := c.Params("uuid")
	id := c.Params("id")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	if err := update(&schema, &table, &id, &uuid, body); err != nil {
		utilities.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(204)
}
