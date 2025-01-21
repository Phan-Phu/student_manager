package schemas

import (
	"student_manager/models/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentResponse struct {
	StudentID primitive.ObjectID `json:"student_id"`
	Name      string             `json:"name"`
	Class     ClassResponse      `json:"class"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateStudent struct {
	StudentID primitive.ObjectID `json:"student_id"`
	Name      string             `json:"name"`
	Class     ClassResponse      `json:"class"`
	BirthDay  string             `json:"birth_day"`
}

type ResponseUpdateScoreStudent struct {
	StudentID   primitive.ObjectID `json:"student_id"`
	StudentName string             `json:"student_name"`
	Score       int                `json:"score"`
}

type StudentDetails struct {
	StudentID primitive.ObjectID `json:"student_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Class     ClassResponse      `json:"class" bson:"classes"`
	BirthDay  string             `json:"birth_day" bson:"birth_day"`
}

func MapStudentDetailToStudentResponse(student *StudentDetails) *StudentResponse {
	return &StudentResponse{
		StudentID: student.StudentID,
		Name:      student.Name,
		BirthDay:  student.BirthDay,
		Class:     student.Class,
	}
}

func MapStudentToStudentResponse(student *db.Student, class *ClassResponse) *StudentResponse {
	return &StudentResponse{
		StudentID: student.ID,
		Name:      student.Name,
		BirthDay:  student.BirthDay,
		Class:     *class,
	}
}
