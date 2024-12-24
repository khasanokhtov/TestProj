package main

import (
	"log"
	"os"
	"test-proj/database"
	"test-proj/middleware"
	"test-proj/routes"

	"github.com/gofiber/fiber/v2"
)

func main(){

	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Ошибка при открытии файла логов:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("Сервер запущен...")

	database.ConnectDatabase()

	app := fiber.New()

	app.Post("/register", routes.Register)
	app.Post("/login", routes.Login)

	// Маршуты через jwt
	api := app.Group("/users", middleware.AuthMiddleware)
	api.Get("/:id/status", routes.GetUserStatus)
	api.Post("/:id/task/complete", routes.CompleteTask)
	api.Post(":id/referrer", routes.AddReferrer)
	api.Get("/leaderboard", routes.GetLeaderboard)

	app.Listen(":3000")

}