package db

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type ClassStatus int

const (
	Opened   ClassStatus = 1
	Closed   ClassStatus = 2
	Full     ClassStatus = 3
	Maintain ClassStatus = 4
)

const MaxStudentInClass = 10

type Class struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	CreateDate       string `bson:"create_date"`
	MaxStudent       int    `bson:"max_student"`
}

type ClassEnrollment struct {
	mgm.DefaultModel `bson:",inline"`
	StudentIds       []string    `bson:"student_ids"`
	ClassId          string      `bson:"class_id"`
	StartDate        time.Time   `bson:"start_date"`
	EndDate          time.Time   `bson:"end_date"`
	Status           ClassStatus `bson:"status"`
}

type OldClassEnrollment struct {
	Data []ClassEnrollment `bson:"data"`
}

func CreateClass(classId int, name string, createDate string) *Class {
	return &Class{
		Name:       name,
		CreateDate: createDate,
	}
}
