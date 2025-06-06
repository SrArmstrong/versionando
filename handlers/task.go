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

func GetTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("userID").(string)

	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Se requiere ID de tarea"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := config.FirestoreClient.Collection("tasks").Doc(taskID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tarea no encontrada"})
	}

	var task models.Task
	if err := doc.DataTo(&task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar tarea"})
	}

	// Verificar que la tarea pertenece al usuario
	if task.UsuarioID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No autorizado"})
	}

	return c.JSON(task)
}

// UpdateTask actualiza una tarea existente
func UpdateTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("userID").(string)

	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Se requiere ID de tarea"})
	}

	var updateData models.TaskUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar que la tarea existe y pertenece al usuario
	docRef := config.FirestoreClient.Collection("tasks").Doc(taskID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tarea no encontrada"})
	}

	var task models.Task
	if err := doc.DataTo(&task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar tarea"})
	}

	if task.UsuarioID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No autorizado"})
	}

	// Preparar actualización
	updates := make(map[string]interface{})
	if updateData.Titulo != nil {
		updates["titulo"] = *updateData.Titulo
	}
	if updateData.Descripcion != nil {
		updates["descripcion"] = *updateData.Descripcion
	}
	if updateData.Completada != nil {
		updates["completada"] = *updateData.Completada
	}
	if updateData.FechaLimite != nil {
		updates["fecha_limite"] = *updateData.FechaLimite
	}

	// Actualizar timestamp
	updates["updated_at"] = time.Now()

	// Aplicar actualización
	_, err = docRef.Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar tarea"})
	}

	return c.JSON(fiber.Map{"message": "Tarea actualizada exitosamente"})
}

// DeleteTask elimina una tarea
func DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("userID").(string)

	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Se requiere ID de tarea"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar que la tarea existe y pertenece al usuario
	docRef := config.FirestoreClient.Collection("tasks").Doc(taskID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tarea no encontrada"})
	}

	var task models.Task
	if err := doc.DataTo(&task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar tarea"})
	}

	if task.UsuarioID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No autorizado"})
	}

	// Eliminar tarea
	_, err = docRef.Delete(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar tarea"})
	}

	return c.JSON(fiber.Map{"message": "Tarea eliminada exitosamente"})
}
