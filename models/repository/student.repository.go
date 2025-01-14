package repository

import (
	"errors"
	"studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentRepositoryInteface interface {
	Create(student *db.Student) error
	FindByID(studentID string) (*db.Student, error)
	Update(student *db.Student) error
	Delete(studentID string) error
	FindAll() ([]*db.Student, error)
	FindByFindOptions(findOptions *options.FindOptions) ([]db.Student, error)
	Count() (int64, error)
	FindByName(name string) (*db.Student, error)
	CreateMany(students []interface{}) error
	DeleteMany(students []interface{}) error
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

func (r *MongoStudentRepository) CreateMany(students []interface{}) error {
	_, err := r.collection.InsertMany(mgm.Ctx(), students)
	if err != nil {
		return errors.New("failed to create student")
	}
	return nil
}

func (r *MongoStudentRepository) DeleteMany(students []interface{}) error {
	_, err := r.collection.DeleteMany(mgm.Ctx(), students)
	if err != nil {
		return errors.New("failed to create student")
	}
	return nil
}

// FindByID retrieves a student by their StudentID
func (r *MongoStudentRepository) FindByID(studentID string) (*db.Student, error) {
	student := &db.Student{}
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}
	err = mgm.Coll(student).First(bson.M{"_id": objID}, student)
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
func (r *MongoStudentRepository) Delete(studentID string) error {
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

func (r *MongoStudentRepository) FindByFindOptions(findOptions *options.FindOptions) ([]db.Student, error) {
	var students []db.Student
	err := mgm.Coll(&db.Student{}).SimpleFind(&students, bson.M{}, findOptions)
	if err != nil {
		return nil, errors.New("cannot find student")
	}
	return students, nil
}

func (r *MongoStudentRepository) Count() (int64, error) {
	total, err := mgm.Coll(&db.Student{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err != nil {
		return 0, errors.New("not Count student")
	}
	return total, nil
}

func (r *MongoStudentRepository) FindByName(name string) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"name": name}, student)
	if err != nil {
		return nil, err
	}
	return student, nil
}
