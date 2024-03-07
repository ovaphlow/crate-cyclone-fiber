package schema

import (
	"ovaphlow/cratecyclone/utility"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetSchema(c *fiber.Ctx) error {
	schemas, err := retrieveSchemas()
	if err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	if len(schemas) > 0 {
		return c.JSON(schemas)
	}
	return c.JSON([]string{})
}

func GetTable(c *fiber.Ctx) error {
	schema := c.Params("schema", "")
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

func Post(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	if err := create(&schema, &table, body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(201)
}

func Get(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	option := c.Query("option", "")
	if option == "default" {
		take, err := strconv.Atoi(c.Query("take", "10"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		order := c.Query("order", "id desc")
		o := QueryOption{Take: take, Skip: int64((page - 1) * take), Order: &order}
		var f QueryFilter
		query := c.Queries()
		if query["equal"] != "" {
			f.Equal = strings.Split(query["equal"], ",")
		}
		if query["not-equal"] != "" {
			f.NotEqual = strings.Split(query["not-equal"], ",")
		}
		if query["like"] != "" {
			f.Like = strings.Split(query["like"], ",")
		}
		if query["greater"] != "" {
			f.Greater = strings.Split(query["greater"], ",")
		}
		if query["lesser"] != "" {
			f.Lesser = strings.Split(query["lesser"], ",")
		}
		if query["in"] != "" {
			f.In = strings.Split(query["in"], ",")
		}
		if query["not-in"] != "" {
			f.NotIn = strings.Split(query["not-in"], ",")
		}
		if query["object-contain"] != "" {
			f.ObjectContain = strings.Split(query["object-contain"], ",")
		}
		if query["array-contain"] != "" {
			f.ArrayContain = strings.Split(query["array-contain"], ",")
		}
		if query["object-like"] != "" {
			f.ObjectLike = strings.Split(query["object-like"], ",")
		}
		result, err := retrieve(&schema, &table, &f, &o)
		if err != nil {
			utility.Slogger.Error(err.Error())
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
	return c.SendString("[]")
}

func GetWithParams(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	uuid := c.Params("uuid")
	id := c.Params("id")
	result, err := retrieveByID(&schema, &table, &id, &uuid)
	if err != nil {
		utility.Slogger.Error(err.Error())
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

func Put(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	uuid := c.Params("uuid")
	id := c.Params("id")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}
	if err := update(&schema, &table, &id, &uuid, body); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(204)
}

func Delete(c *fiber.Ctx) error {
	schema := c.Params("schema")
	table := c.Params("table")
	uuid := c.Params("uuid")
	id := c.Params("id")
	if err := remove(&schema, &table, &id, &uuid); err != nil {
		utility.Slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(204)
}
