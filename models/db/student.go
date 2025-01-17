package db

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string             `bson:"name"`
	Age              int                `bson:"age"`
	BirthDay         string             `bson:"birth_day"`
	Score            int                `bson:"score"`
	ClassId          primitive.ObjectID `bson:"class_id"`
}

func NewStudent(name string, classId primitive.ObjectID, birthDay string, age int, score int) *Student {
	return &Student{
		Name:     name,
		Age:      age,
		BirthDay: birthDay,
		Score:    score,
		ClassId:  classId,
	}
}
