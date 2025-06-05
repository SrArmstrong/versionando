package routes

import (
	"versionando/handlers"
	"versionando/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Rutas públicas
	api := app.Group("/api")
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	// Rutas protegidas con JWT
	auth := api.Group("/", middleware.JWTProtected())
	auth.Post("/tasks", handlers.CreateTask)
	auth.Get("/tasks", handlers.GetTasks)
}
