// pkg/database/database.go
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/66emil/fiber-todo-api/pkg/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// 2. Формируем строку DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// 3. Открываем соединение GORM
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	log.Println("Успешное подключение к базе данных!")

	// 4. Автоматическая миграция (создание таблиц User и Todo)
	DB.AutoMigrate(&models.User{}, &models.Todo{})
	log.Println("Миграция таблиц завершена.")
}
