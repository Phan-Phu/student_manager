package repository

import (
	"errors"
	"studenent_manager/models"
	"studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepository interface {
	// Create(Class *db.Class) error
	FindByID(classIDString string) (*db.Class, error)
	// Update(Class *db.Class) error
	// Delete(ClassID int) error
	// FindAll() ([]*db.Class, error)
}

type MongoClassRepository struct {
	classCollection *mgm.Collection
}

func NewMongoClassRepository() *MongoClassRepository {
	return &MongoClassRepository{
		classCollection: mgm.Coll(&db.Class{}),
	}
}

// func (r *MongoRepository) Create(Class *db.Class) error {
// 	if err := r.classCollection.Create(Class); err != nil {
// 		return errors.New("failed to create Class")
// 	}
// 	return nil
// }

// // FindByID retrieves a Class by their ClassID
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

// // Update modifies an existing Class's data
// func (r *MongoRepository) Update(Class *db.Class) error {
// 	if err := r.classCollection.Update(Class); err != nil {
// 		return errors.New("failed to update Class")
// 	}
// 	return nil
// }

// // Delete removes a Class from the database
// func (r *MongoRepository) Delete(ClassID int) error {
// 	Class, err := r.FindByID(ClassID)
// 	if err != nil {
// 		return err
// 	}
// 	if err := r.classCollection.Delete(Class); err != nil {
// 		return errors.New("failed to delete Class")
// 	}
// 	return nil
// }

// func (r *MongoRepository) FindAll() ([]*db.Class, error) {
// 	var class []*db.Class
// 	if err := r.classCollection.SimpleFind(&class, bson.M{}); err != nil {
// 		return nil, errors.New("failed to retrieve class")
// 	}
// 	return class, nil
// }

func (r *MongoClassRepository) FindClassByName(name string) (*db.Class, error) {
	Class := &db.Class{}
	err := r.classCollection.First(bson.M{"name": name}, Class)
	if err != nil {
		return nil, err
	}
	return Class, nil
}
