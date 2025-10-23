package handlers

import (
	"os"
	"time"

	"github.com/66emil/fiber-todo-api/pkg/database"
	"github.com/66emil/fiber-todo-api/pkg/models"
	"github.com/66emil/fiber-todo-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthBody struct {
	Email    string `json:"email", xml:"email", form:"email"`
	Password string `json:"password", xml:"password", form:"password"`
}

func Register(c *fiber.Ctx) error {
	body := new(AuthBody)

	// Парсинг JSON тела запроса
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	// Хэшируем пароль
	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка хеширования пароля"})
	}

	// Создаем объект пользователя
	user := models.User{
		Email:        body.Email,
		PasswordHash: hash,
	}

	// Сохраняем в БД
	if result := database.DB.Create(&user); result.Error != nil {
		// Обычно это ошибка, если Email уже существует (Unique Constraint)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Пользователь с таким email уже существует"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Пользователь успешно зарегистрирован", "user_id": user.ID})
}

func Login(c *fiber.Ctx) error {
	body := new(AuthBody)

	// Парсинг JSON тела запроса
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}
	var user models.User

	if result := database.DB.Where("email = ?", body.Email).First(&user); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный email или пароль"})
	}

	if !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный email или пароль"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(), // Токен истекает через 12 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать токен"})
	}

	// 4. Отправляем токен клиенту
	return c.JSON(fiber.Map{"token": t})
}
