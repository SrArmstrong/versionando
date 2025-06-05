package handlers

import (
	"context"
	"time"
	"versionando/config"
	"versionando/models"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Obtener el ID del usuario desde el JWT
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar que el usuario existe
	_, err := config.FirestoreClient.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Usuario no válido"})
	}

	// Asignar fechas y usuario
	task.UsuarioID = userID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Crear tarea
	docRef, _, err := config.FirestoreClient.Collection("tasks").Add(ctx, task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear la tarea"})
	}

	// Actualizar el ID del documento
	_, err = docRef.Set(ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar ID de la tarea"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tarea creada exitosamente",
		"id":      docRef.ID,
	})
}

func GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tasks []models.Task
	iter := config.FirestoreClient.Collection("tasks").Where("usuario_id", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener tareas"})
		}

		var task models.Task
		if err := doc.DataTo(&task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}
