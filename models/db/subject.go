package db

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Subject struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	CreateDate       string `bson:"create_date"`
}

func NewSubject(name string) *Subject {
	return &Subject{
		Name:       name,
		CreateDate: time.Now().String(),
	}
}
