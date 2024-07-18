package controllers

import (
	"github.com/gofiber/fiber/v2"
)


func Test(c *fiber.Ctx) error {
	return c.SendString("Working hehehe 2.123")
}

func AuthHandler(c *fiber.Ctx) error {
	username := c.Locals("username")
	c.SendStatus(200)
	return c.JSON(fiber.Map{"username": username})
}
