package models

import "time"

type Task struct {
	ID          string    `firestore:"id" json:"id"`
	Titulo      string    `firestore:"titulo" json:"titulo"`
	Descripcion string    `firestore:"descripcion" json:"descripcion"`
	FechaInicio time.Time `firestore:"fecha_inicio" json:"fecha_inicio"`
	FechaLimite time.Time `firestore:"fecha_limite" json:"fecha_limite"`
	UsuarioID   string    `firestore:"usuario_id" json:"usuario_id"`
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updated_at"`
}

type TaskUpdate struct {
	Titulo      *string    `json:"titulo,omitempty" firestore:"titulo,omitempty"`
	Descripcion *string    `json:"descripcion,omitempty" firestore:"descripcion,omitempty"`
	Completada  *bool      `json:"completada,omitempty" firestore:"completada,omitempty"`
	FechaLimite *time.Time `json:"fecha_limite,omitempty" firestore:"fecha_limite,omitempty"`
}
