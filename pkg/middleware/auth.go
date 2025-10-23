package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is missing"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		// Ошибка 401: Неверный формат заголовка
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный формат заголовка Authorization"})
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Убеждаемся, что используется ожидаемый алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный метод подписи токена")
		}
		// Возвращаем секретный ключ для валидации
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный или истекший токен"})
	}

	// Извлекаем пользовательский ID из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Ошибка при получении данных токена"})
	}

	if userID, ok := claims["user_id"].(float64); ok {
		c.Locals("user_id", uint(userID))
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Отсутствует ID пользователя в токене"})
	}

	// 5. Передаем управление следующему хэндлеру/мидлварю
	return c.Next()
}
