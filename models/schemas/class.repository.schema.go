package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type ClassResponse struct {
	ClassID    primitive.ObjectID `bson:"_id" json:"class_id"`
	CreateDate string             `bson:"create_date" json:"create_date"`
	Name       string             `bson:"name" json:"name"`
	MaxStudent int                `bson:"max_student" json:"max_student"`
}

type RequestClass struct {
	Name       string `json:"name"`
	MaxStudent int    `json:"max_student"`
}

type UpdateClassRequest struct {
	ID         string `json:"class_id"`
	Name       string `json:"name"`
	MaxStudent int    `json:"max_student"`
}

type ResponseClassUpdate struct {
	ID         string `json:"class_id"`
	Name       string `json:"name"`
	MaxStudent int    `json:"max_student"`
}

type DeleteClassRequest struct {
	ID string `json:"class_id"`
}
