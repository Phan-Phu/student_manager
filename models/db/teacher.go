package db

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role int

const (
	AdminRole   Role = 0 // admin can be fix in first release.
	TeacherRole Role = 1
)

type Teacher struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string               `bson:"name"`
	BirthDay         string               `bson:"birth_day"`
	Username         string               `bson:"username"`
	Password         string               `bson:"password"`
	ClassIds         []primitive.ObjectID `bson:"class_ids"`
	SubjectIds       []primitive.ObjectID `bson:"subject_ids"`
}

func NewTeacher(name string, classIds []primitive.ObjectID, birthDay string, userName string, password string) *Teacher {
	return &Teacher{
		Name:     name,
		ClassIds: classIds,
		BirthDay: birthDay,
		Username: userName,
		Password: password,
	}
}
