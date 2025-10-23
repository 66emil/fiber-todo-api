package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetTodos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	return c.JSON(fiber.Map{
		"message": "Получение списка дел для пользователя",
		"user_id": userID,
		"data":    "[]",
	})
}
