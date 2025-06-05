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

	// Rutas para CRUD
	auth.Post("/tasks", handlers.CreateTask)
	auth.Get("/tasks", handlers.GetTasks)

	// Rutas CRUD para usuarios
	auth.Get("/users", handlers.GetAllUsers)
	auth.Get("/users/:id", handlers.GetUser)
	auth.Put("/users/:id", handlers.UpdateUser)
	auth.Delete("/users/:id", handlers.DeleteUser)
}
