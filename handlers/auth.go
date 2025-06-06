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

func GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Se requiere ID de usuario",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := config.FirestoreClient.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al procesar datos del usuario",
		})
	}

	// No devolver la contraseña ni la respuesta secreta
	user.Password = ""
	user.RespuestaSecreta = ""

	return c.JSON(user)
}

// UpdateUser actualiza un usuario existente
func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Se requiere ID de usuario",
		})
	}

	var updateData models.UserUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar que el usuario existe
	docRef := config.FirestoreClient.Collection("users").Doc(userID)
	if _, err := docRef.Get(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Preparar datos para actualizar
	updates := make(map[string]interface{})
	if updateData.Nombre != "" {
		updates["nombre"] = updateData.Nombre
	}
	if updateData.Apellidos != "" {
		updates["apellidos"] = updateData.Apellidos
	}
	if updateData.Email != "" {
		// Verificar si el nuevo email ya existe
		iter := config.FirestoreClient.Collection("users").
			Where("email", "==", updateData.Email).
			Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error al verificar email",
				})
			}
			if doc.Ref.ID != userID {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "El email ya está en uso por otro usuario",
				})
			}
		}
		updates["email"] = updateData.Email
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al encriptar la contraseña",
			})
		}
		updates["password"] = string(hashedPassword)
	}

	// Actualizar timestamp
	updates["updatedAt"] = time.Now()

	// Realizar la actualización
	_, err := docRef.Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el usuario",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Usuario actualizado exitosamente",
	})
}

// DeleteUser elimina un usuario
func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Se requiere ID de usuario",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar que el usuario existe
	docRef := config.FirestoreClient.Collection("users").Doc(userID)
	_, err := docRef.Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Eliminar el usuario
	_, err = docRef.Delete(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar el usuario",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Usuario eliminado exitosamente",
	})
}

// GetAllUsers obtiene todos los usuarios (solo para administradores)
func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iter := config.FirestoreClient.Collection("users").Documents(ctx)
	var users []models.UserPublic

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al obtener usuarios",
			})
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al procesar datos del usuario",
			})
		}

		// Crear versión pública del usuario sin información sensible
		publicUser := models.UserPublic{
			ID:        user.ID,
			Nombre:    user.Nombre,
			Apellidos: user.Apellidos,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		users = append(users, publicUser)
	}

	return c.JSON(users)
}

func GetSecretQuestion(c *fiber.Ctx) error {
	var request models.PasswordRecoveryRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar usuario por email
	iter := config.FirestoreClient.Collection("users").Where("email", "==", request.Email).Limit(1).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil || len(docs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error procesando datos del usuario",
		})
	}

	return c.JSON(fiber.Map{
		"preguntaSecreta": user.PreguntaSecreta,
	})
}

// RecoverPassword permite actualizar la contraseña si la respuesta secreta es correcta
func RecoverPassword(c *fiber.Ctx) error {
	var input models.SecretAnswerVerification
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar usuario por email
	iter := config.FirestoreClient.Collection("users").Where("email", "==", input.Email).Limit(1).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil || len(docs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error procesando datos del usuario",
		})
	}

	// Verificar respuesta secreta
	if err := bcrypt.CompareHashAndPassword([]byte(user.RespuestaSecreta), []byte(input.RespuestaSecreta)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Respuesta secreta incorrecta",
		})
	}

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NuevaPassword), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo hashear la nueva contraseña",
		})
	}

	// Actualizar la contraseña en la base de datos
	_, err = docs[0].Ref.Update(ctx, []firestore.Update{
		{Path: "password", Value: string(hashedPassword)},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo actualizar la contraseña",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Contraseña actualizada exitosamente",
	})
}
