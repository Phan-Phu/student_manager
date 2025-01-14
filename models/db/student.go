package db

import (
	"github.com/kamva/mgm/v3"
)

type Student struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	Age              int    `bson:"age"`
	ClassID          int    `bson:"class_id"`
	BirthDay         string `bson:"birth_day"`
	Score            int    `bson:"score"`
}

func NewStudent(name string, classID int, birthDay string, age int, score int) *Student {
	return &Student{
		Name:     name,
		ClassID:  classID,
		Age:      age,
		BirthDay: birthDay,
		Score:    score,
	}
}
