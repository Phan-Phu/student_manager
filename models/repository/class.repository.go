package repository

import (
	"errors"
	"student_manager/models"
	"student_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepositoryInterface interface {
	Create(class *db.Class) error
	Update(class *db.Class) error
	FindByID(ClassId string) (*db.Class, error)
	FindAll() ([]*db.Class, error)
}

type ClassRepository struct {
	collection *mgm.Collection
}

func NewMongoClassRepository() *ClassRepository {
	return &ClassRepository{
		collection: mgm.Coll(&db.Class{}),
	}
}

func (r *ClassRepository) Create(class *db.Class) error {
	return r.collection.Create(class)
}

func (r *ClassRepository) Update(class *db.Class) error {
	return r.collection.Update(class)
}

func (r *ClassRepository) FindByID(classId primitive.ObjectID) (*db.Class, error) {
	Class := &db.Class{}
	err := r.collection.FindByID(classId, Class)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeNotFoundClass))
	}
	return Class, nil
}

func (r *ClassRepository) FindAll() ([]*db.Class, error) {
	var classes []*db.Class
	if err := r.collection.SimpleFind(&classes, bson.M{}); err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}
	return classes, nil
}
