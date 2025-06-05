package main

import (
	"versionando/config"
	"versionando/handlers"
	"versionando/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Conectar a la base de datos
	config.ConnectDB()

	// Rutas públicas
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Rutas protegidas
	api := app.Group("/api", middleware.JWTProtected())
	api.Post("/tasks", handlers.CreateTask)
	// Agrega aquí las demás rutas protegidas (CRUD de tareas y usuarios)

	app.Listen(":3000")
}
