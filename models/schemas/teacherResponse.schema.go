package schemas

import (
	"student_manager/models/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherResponse struct {
	TeacherID primitive.ObjectID `json:"teacher_id"`
	Name      string             `json:"name"`
	Classes   []*ClassResponse   `json:"classes"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateTeacher struct {
	TeacherID primitive.ObjectID `json:"teacher_id"`
	Name      string             `json:"name"`
	Classes   []*ClassResponse   `json:"classes"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateScoreTeacher struct {
	TeacherID   primitive.ObjectID `json:"teacher_id"`
	TeacherName string             `json:"name"`
	Score       int                `json:"score"`
}

type TeacherDetails struct {
	TeacherID primitive.ObjectID `json:"teacher_id"`
	Name      string             `json:"name"`
	Classes   []*ClassResponse   `json:"classes"`
	BirthDay  string             `json:"birth_day"`
}

func MapTeacherDetailToTeacherResponse(Teacher *TeacherDetails) *TeacherResponse {
	return &TeacherResponse{
		TeacherID: Teacher.TeacherID,
		Name:      Teacher.Name,
		BirthDay:  Teacher.BirthDay,
		Classes:   Teacher.Classes,
	}
}

func MapTeacherToTeacherResponse(Teacher *db.Teacher, classes []*ClassResponse) *TeacherResponse {
	return &TeacherResponse{
		TeacherID: Teacher.ID,
		Name:      Teacher.Name,
		BirthDay:  Teacher.BirthDay,
		Classes:   classes,
	}
}
