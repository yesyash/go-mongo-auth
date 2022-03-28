package gomongoauth

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name     string             `json:"name" validate:"required,min=2,max=100"`
	Email    string             `json:"email" validate:"required,email"`
	Password string             `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Err    error  `json:"err"`
	Data   string `json:"data,omitempty"`
}
