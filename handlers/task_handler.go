package handlers

import (
	"context"
	"time"
	"versionando/config"
	"versionando/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Obtener el ID del usuario desde el contexto (establecido por el middleware JWT)
	userID := c.Locals("userID").(string)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de usuario inválido"})
	}
	task.UsuarioID = oid

	collection := config.DB.Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear la tarea"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Tarea creada exitosamente"})
}
