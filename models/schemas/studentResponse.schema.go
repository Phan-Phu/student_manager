package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResponseStudent struct {
	StudentID primitive.ObjectID `json:"student_id"`
	Name      string             `json:"name"`
	ClassID   int                `json:"class_id"`
	Age       int                `json:"age"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateStudent struct {
	StudentID primitive.ObjectID `json:"student_id"`
	Name      string             `json:"name"`
	ClassID   int                `json:"class_id"`
	Age       int                `json:"age"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateScoreStudent struct {
	StudentID   primitive.ObjectID `json:"student_id"`
	StudentName string             `json:"student_name"`
	Score       int                `json:"score"`
}
