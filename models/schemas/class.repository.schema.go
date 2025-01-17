package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type ClassResponse struct {
	ClassID    primitive.ObjectID `json:"class_id"`
	CreateDate string             `json:"create_date"`
	Name       string             `json:"name"`
	MaxStudent int                `json:"max_student"`
}

type ClassRequest struct {
	Name       string `json:"name"`
	MaxStudent int    `json:"max_student"`
}
