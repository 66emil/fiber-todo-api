package handlers

import (
	"github.com/66emil/fiber-todo-api/pkg/database"
	"github.com/66emil/fiber-todo-api/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type TodoBody struct {
	Title  string `json"title", xml:"title", form:"title"`
	IsDone *bool  `json"is_done", xml:"is_done", form:"is_done"`
}

func CreateTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	body := new(TodoBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	todo := models.Todo{
		UserID: userID,
		Title:  body.Title,
		IsDone: false,
	}

	if result := database.DB.Create(&todo); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при создании задачи"})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func GetTodos(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var todos []models.Todo

	database.DB.Where("user_id = ?", userID).Find(&todos)

	return c.JSON(todos)
}

func UpdateTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todoID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Нет валидного ID задачи"})
	}

	body := new(TodoBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	var todo models.Todo

	if result := database.DB.First(&todo, "id = ? AND user_id = ?", todoID, userID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задача не найдена"})
	}

	if body.Title != "" {
		todo.Title = body.Title
	}
	if body.IsDone != nil {
		todo.IsDone = *body.IsDone
	}

	database.DB.Save(&todo)

	return c.JSON(todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todoID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID задачи"})
	}

	var todo models.Todo

	// 1. Проверяем наличие и права доступа
	if result := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задача не найдена или вам недоступна"})
	}

	// 2. Удаляем задачу
	database.DB.Delete(&todo)

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content - успешное удаление без тела ответа
}
