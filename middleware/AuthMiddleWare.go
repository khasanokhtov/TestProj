package middleware

import (
	"log"
	"test-proj/utils"

	"strings"

	"github.com/gofiber/fiber/v2"
)

// Валидация токена

func AuthMiddleware(c *fiber.Ctx) error {
	//извлекаем из загаловка токен
	
	tokenString := c.Get("Authorization")
	if tokenString == ""{
		log.Println("Ошибка: токен отсутствует")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "токен отсутствует"})
	}
	log.Println("Полученный токен:", tokenString)
	// удаляем префикс "Bearer "
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		log.Println("Ошибка: токен отсутствует после удаления префикса")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "токен отсутствует"})
	}

	//Проверяем
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		log.Println("Ошибка: недействительный токен")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "недействительный токен"})
	}

	log.Println("Успешно проверен токен для userID:", claims.UserID)

	c.Locals("UserID", float64(claims.UserID))
	log.Println("Сохранённый в контекст userID:", c.Locals("userID"))
	return c.Next()
}