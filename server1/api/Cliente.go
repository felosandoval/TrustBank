package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cliente struct {
	ID                    primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Nombre                string             `bson:"nombre,omitempty" json:"nombre,omitempty"`
	Contrasena            string             `bson:"contrasena,omitempty" json:"contrasena,omitempty"`
	Fecha_nacimiento      string             `bson:"fecha_nacimiento,omitempty" json:"fecha_nacimiento,omitempty"`
	Direccion             string             `bson:"direccion,omitempty" json:"direccion,omitempty"`
	Numero_identificacion string             `bson:"numero_identificacion,omitempty" json:"numero_identificacion,omitempty"`
	Email                 string             `bson:"email,omitempty" json:"email,omitempty"`
	Telefono              string             `bson:"telefono,omitempty" json:"telefono,omitempty"`
	Genero                string             `bson:"genero,omitempty" json:"genero,omitempty"`
	Nacionalidad          string             `bson:"nacionalidad,omitempty" json:"nacionalidad,omitempty"`
	Ocupacion             string             `bson:"ocupacion,omitempty" json:"ocupacion,omitempty"`
}
