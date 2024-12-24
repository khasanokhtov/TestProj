package routes

import (
	"log"
	"test-proj/database"
	"test-proj/models"
	"test-proj/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

//регистрация
func Register(c *fiber.Ctx) error {
	// получаем запрос
	var input struct {
		Name 	 string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный запрос"})
	}

	//проверка существования пользователя
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "пользователь с таким email уже существует"})
	}

	//Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось хэшировать пароль"})
	}

	//Создаем нового пользователя
	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: string(hashedPassword),
		Balance: 0,	
	}

	//Сохраняем в БД
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось создать токен"})
	}

	//Генерируем JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось создать токен"})
	}

	return c.JSON(fiber.Map{"token": token, "user": user})
}

//Логин
func Login(c *fiber.Ctx) error {
	// Данные запроса
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"` 
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный запрос"})
	}

	//Поиск пользователя по базе
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не найден"})
	}

	//проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "неверный пароль"})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось создать токен"})
	}

	return c.JSON(fiber.Map{"token": token})
}

//Информация о пользователе
func GetUserStatus(c *fiber.Ctx) error{
	
	log.Println("Запрос информации о пользователе")

	//извлекаем userID
	log.Println("Контекст перед извлечением userID:", c.Locals("UserID"))
	userIDInterface  := c.Locals("UserID")
	if userIDInterface == nil {
		log.Println("Ошибка: пользователь не авторизован")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не авторизован"})
	}

	// преобразуем интерфейс в uint
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		log.Println("Ошибка: неверный ID пользователя")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "неверный ID пользователя"})
	}
	userID := uint(userIDFloat)

	// Получаем инф из БД
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Println("Ошибка: пользователь не найден")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "пользователь не найден"})
	}

	log.Println("Информация о пользователе получена:", user)

	return c.JSON(user)
}

//Вполнение задания
func CompleteTask(c *fiber.Ctx) error {
	//опять извлекаем пользователя 
	userID := c.Locals("UserID")

	//Получаем из БД пользователя
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "пользователь не найдет"})
	}

	// Добавляем баллы за выполнения задания
	user.Balance += 10

	//Сохраняем
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось обновить баланс"})
	}

	return c.JSON(fiber.Map{"message": "задание выполнено", "Баланс": user.Balance})
}

func AddReferrer(c *fiber.Ctx) error {
	// Извлекаем userID
	userIDInterface := c.Locals("UserID")
	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не авторизован"})
	}

	userIDFloat, ok := userIDInterface.(float64) // Проверяем тип float64
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "неверный ID пользователя"})
	}

	
	userID := uint(userIDFloat)
	

	// Получаем даннные запроса
	var input struct {
		ReferrerID uint `json:"referrer_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный запрос"})
	}

	// Получаем пользователя из базы данных
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "пользователь не найден"})
	}

	// Проверяем, что у пользователя нет реферального ID
	if user.ReferrerID != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "реферальный код уже добавлен"})
	}

	// Проверяем, существует ли реферальный пользователь
	var referrer models.User
	if err := database.DB.First(&referrer, input.ReferrerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "реферальный пользователь не найден"})
	}

	// Добавляем реферальный ID и обновляем баланс
	user.ReferrerID = &input.ReferrerID
	user.Balance += 5 

	// Сохраняем изменения
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось добавить реферальный код"})
	}

	return c.JSON(fiber.Map{"message": "реферальный код добавлен", "balance": user.Balance})
}


//Топ пользователей 
func GetLeaderboard(c * fiber.Ctx) error {
	var users []models.User

	//Сортируем
	if err := database.DB.Order("balance desc").Limit(5).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "не удалось получать рейтинг"})
	}

	return c.JSON(users)
}