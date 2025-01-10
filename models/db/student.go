package db

import "github.com/kamva/mgm/v3"

type Student struct {
	mgm.DefaultModel `bson:",inline"`
	StudentID        int    `json:"student_id" bson:"student_id"`
	Name             string `json:"name" bson:"name"`
	ClassID          int    `json:"class_id" bson:"class_id"`
	BirthDay         string `json:"birth_day" bson:"birth_day"`
	Score            int    `json:"score" bson:"score"`
	Role             Role   `json:"role" bson:"role"` // not input from user
}

func NewStudent(studentId int, name string, classID int, birthDay string, score int) *Student {
	return &Student{
		//ID:        primitive.NewObjectID(),
		StudentID: studentId,
		Name:      name,
		ClassID:   classID,
		BirthDay:  birthDay,
		Score:     score,
		Role:      StudentRole,
	}
}
