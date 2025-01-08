package db

import "github.com/kamva/mgm/v3"

type Role int

const (
	NoneRole    Role = -1
	AdminRole   Role = 0 // admin can be fix in first release.
	TeacherRole Role = 1
	StudentRole Role = 2
)

type Teacher struct {
	mgm.DefaultModel `bson:",inline"`
	TeacherID        int    `json:"teacher_id" bson:"teacher_id"`
	Name             string `json:"name" bson:"name"`
	BirthDay         string `json:"birth_day" bson:"birth_day"`
	ClassIds         []int  `json:"class" bson:"class"`
	Role             Role   `json:"role" bson:"role"` // not input from user
	Username         string `json:"username" bson:"username"`
	Password         string `json:"password" bson:"password"`
}

func CreateTeacher(teacherId int, name string, classIds []int, birthDay string, userName string, password string) *Teacher {
	return &Teacher{

		TeacherID: teacherId,
		Name:      name,
		ClassIds:  classIds,
		BirthDay:  birthDay,
		Role:      TeacherRole,
		Username:  userName,
		Password:  password,
	}
}
