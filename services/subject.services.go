package services

import (
	"fmt"
	"student_manager/models/db"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
)

func CreateSubject(data schemas.SubjectRequest) (*schemas.SubjectResponse, error) {
	subject := db.NewSubject(data.Name)

	err := mgm.Coll(subject).Create(subject)
	if err != nil {
		return nil, fmt.Errorf("failed to save subject: %v", err)
	}

	response := &schemas.SubjectResponse{
		ID:         subject.ID.Hex(),
		Name:       subject.Name,
		CreateDate: subject.CreateDate,
	}

	return response, nil
}
