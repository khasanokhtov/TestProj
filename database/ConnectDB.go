package database

import (
	"fmt"
	"log"
	"os"
	"test-proj/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Глобальная для подключения к базе
var DB *gorm.DB

// Функция подключения к постгре
func ConnectDatabase(){
	//переменные окружения .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	//Получения поараметров подключения из переменных окружения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Коннект к базе данныъс
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе: ", err)
	}

	DB = database
	fmt.Println("Подключение к базе данных успешно")

	//Автомиграция
	database.AutoMigrate(&models.User{})
}