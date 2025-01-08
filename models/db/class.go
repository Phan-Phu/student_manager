package db

import "github.com/kamva/mgm/v3"

type Class struct {
	mgm.DefaultModel `bson:",inline"`
	ClassID          int       `json:"id" bson:"class_id"`
	Name             string    `json:"name" bson:"name"`
	Students         []Student `bson:"-"`          // not need store in database
	Teachers         []Teacher `bson:"-"`          // not need store in database
	StudentIds       []int     `bson:"studentIds"` // not need store in database
	TeacherIds       []int     `bson:"teacherIds"` // not need store in database
	CreateDate       string    `json:"create_date" bson:"create_date"`
	Role             Role      `json:"role" bson:"role"`
}

func CreateClass(classId int,
	name string,
	createDate string,
	studentIds []int,
	teacherIds []int,
) *Class {
	return &Class{
		ClassID:    classId,
		Name:       name,
		CreateDate: createDate,
		Students:   nil,
		Teachers:   nil,
		StudentIds: studentIds,
		TeacherIds: teacherIds,
		Role:       NoneRole,
	}
}
