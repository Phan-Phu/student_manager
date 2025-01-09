package repository

import (
	"errors"
	"studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type StudentRepository interface {
	Create(student *db.Student) error
	FindByID(studentID int) (*db.Student, error)
	Update(student *db.Student) error
	Delete(studentID int) error
	FindAll() ([]*db.Student, error)
}

type MongoStudentRepository struct {
	collection *mgm.Collection
}

func NewMongoStudentRepository() *MongoStudentRepository {
	return &MongoStudentRepository{
		collection: mgm.Coll(&db.Student{}),
	}
}

func (r *MongoStudentRepository) Create(student *db.Student) error {
	if err := r.collection.Create(student); err != nil {
		return errors.New("failed to create student")
	}
	return nil
}

// FindByID retrieves a student by their StudentID
func (r *MongoStudentRepository) FindByID(studentID int) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"student_id": studentID}, student)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// Update modifies an existing student's data
func (r *MongoStudentRepository) Update(student *db.Student) error {
	if err := r.collection.Update(student); err != nil {
		return errors.New("failed to update student")
	}
	return nil
}

// Delete removes a student from the database
func (r *MongoStudentRepository) Delete(studentID int) error {
	student, err := r.FindByID(studentID)
	if err != nil {
		return err
	}
	if err := r.collection.Delete(student); err != nil {
		return errors.New("failed to delete student")
	}
	return nil
}

func (r *MongoStudentRepository) FindAll() ([]*db.Student, error) {
	var students []*db.Student
	if err := r.collection.SimpleFind(&students, bson.M{}); err != nil {
		return nil, errors.New("failed to retrieve students")
	}
	return students, nil
}
