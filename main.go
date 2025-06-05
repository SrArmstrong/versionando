package main

import (
	"log"
	"versionando/config"
	"versionando/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Inicializar Firebase
	config.ConnectFirestore()
	defer config.CloseFirestore()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())

	// Configurar rutas
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
