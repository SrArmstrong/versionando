package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Titulo      string             `bson:"titulo" json:"titulo"`
	Descripcion string             `bson:"descripcion" json:"descripcion"`
	FechaInicio string             `bson:"fecha_inicio" json:"fecha_inicio"`
	FechaLimite string             `bson:"fecha_limite" json:"fecha_limite"`
	UsuarioID   primitive.ObjectID `bson:"usuario_id" json:"usuario_id"`
}
