package handlers

import (
	"context"
	"time"
	"versionando/config"
	"versionando/models"
	"versionando/utils"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func Register(c *fiber.Ctx) error {
	var user models.User

	// Parsear el JSON del cuerpo de la petición al struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	// Verificar si ya existe un usuario con el mismo correo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iter := config.FirestoreClient.Collection("users").Where("email", "==", user.Email).Documents(ctx)
	if _, err := iter.Next(); err != iterator.Done {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "El email ya está registrado",
		})
	}

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al encriptar la contraseña",
		})
	}
	user.Password = string(hashedPassword)

	// Encriptar respuesta secreta
	hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(user.RespuestaSecreta), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al encriptar la respuesta secreta",
		})
	}
	user.RespuestaSecreta = string(hashedAnswer)

	// Asignar timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insertar usuario
	docRef, _, err := config.FirestoreClient.Collection("users").Add(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al registrar el usuario",
		})
	}

	// Asignar ID al documento creado
	_, err = docRef.Set(ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el ID del usuario",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuario registrado exitosamente",
		"id":      docRef.ID,
	})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iter := config.FirestoreClient.Collection("users").Where("email", "==", input.Email).Limit(1).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil || len(docs) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error procesando datos del usuario"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
