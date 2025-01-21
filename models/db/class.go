package db

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type ClassStatus int

const (
	ClassStatusOpened   ClassStatus = 1
	ClassStatusClosed   ClassStatus = 2
	ClassStatusFull     ClassStatus = 3
	ClassStatusMaintain ClassStatus = 4
)

const MaxStudentInClass = 10

type Class struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	CreateDate       string `bson:"create_date"`
	MaxStudent       int    `bson:"current_student"`
}

func NewClass(name string, maxStudent int) *Class {
	return &Class{
		Name:       name,
		CreateDate: time.November.String(),
		MaxStudent: maxStudent,
	}
}
