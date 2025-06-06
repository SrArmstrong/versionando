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
	api.Post("/passrecovery", handlers.RecoverPassword)

	// Nuevas rutas para recuperación de contraseña
	api.Post("/passrecovery/question", handlers.GetSecretQuestion)

	// Rutas protegidas con JWT
	auth := api.Group("/", middleware.JWTProtected())

	// Rutas para CRUD
	auth.Post("/tasks", handlers.CreateTask)
	auth.Get("/tasks", handlers.GetTasks)
	auth.Get("/tasks/:id", handlers.GetTask)
	auth.Put("/tasks/:id", handlers.UpdateTask)
	auth.Delete("/tasks/:id", handlers.DeleteTask)

	// Rutas CRUD para usuarios
	auth.Get("/users", handlers.GetAllUsers)
	auth.Get("/users/:id", handlers.GetUser)
	auth.Put("/users/:id", handlers.UpdateUser)
	auth.Delete("/users/:id", handlers.DeleteUser)
}
