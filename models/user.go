package models

import "time"

type User struct {
	ID               string    `firestore:"id" json:"id"`
	Nombre           string    `firestore:"nombre" json:"nombre"`
	Apellidos        string    `firestore:"apellidos" json:"apellidos"`
	Email            string    `firestore:"email" json:"email"`
	Password         string    `firestore:"password" json:"password"`
	FechaNacimiento  time.Time `firestore:"fecha_nacimiento" json:"fecha_nacimiento"`
	PreguntaSecreta  string    `firestore:"pregunta_secreta" json:"pregunta_secreta"`
	RespuestaSecreta string    `firestore:"respuesta_secreta" json:"respuesta_secreta"`
	CreatedAt        time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt        time.Time `firestore:"updated_at" json:"updated_at"`
}

type UserUpdate struct {
	Nombre    string `json:"nombre,omitempty"`
	Apellidos string `json:"apellidos,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type UserPublic struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
