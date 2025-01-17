package repository

import (
	"errors"
	"studenent_manager/models"
	"studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepository interface {
	FindByID(classIDString string) (*db.Class, error)
}

type MongoClassRepository struct {
	classCollection *mgm.Collection
}

func NewMongoClassRepository() *MongoClassRepository {
	return &MongoClassRepository{
		classCollection: mgm.Coll(&db.Class{}),
	}
}

func (r *MongoClassRepository) FindByID(classIDString string) (*db.Class, error) {
	classId, err := primitive.ObjectIDFromHex(classIDString)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	Class := &db.Class{}
	err = r.classCollection.FindByID(classId, Class)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeNotFoundClass))
	}
	return Class, nil
}
