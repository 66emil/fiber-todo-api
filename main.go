// main.go
package main

import (
	"log"

	"github.com/66emil/fiber-todo-api/pkg/database"
	"github.com/66emil/fiber-todo-api/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	// Группа маршрутов для аутентификации
	authGroup := app.Group("/auth")
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)

	// Здесь позже добавим защищенные маршруты
}

func main() {
	// 1. Подключаемся к базе данных и выполняем миграцию
	database.ConnectDB()

	app := fiber.New()

	setupRoutes(app)

	// 2. Запуск сервера с обработкой ошибок
	log.Println("Сервер запускается на порту :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
