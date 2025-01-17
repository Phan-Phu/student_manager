package repository

// import (
// 	"errors"
// 	"studenent_manager/models/db"

// 	"github.com/kamva/mgm/v3"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// type TeacherRepository interface {
// 	Create(teacher *db.Teacher) error
// 	FindByID(teacherID int) (*db.Teacher, error)
// 	Update(teacher *db.Teacher) error
// 	Delete(teacherID int) error
// 	FindAll() ([]*db.Teacher, error)
// }

// type MongoTeacherRepository struct {
// 	collection *mgm.Collection
// }

// func NewMongoTeacherRepository() *MongoTeacherRepository {
// 	return &MongoTeacherRepository{
// 		collection: mgm.Coll(&db.Teacher{}),
// 	}
// }

// func (r *MongoTeacherRepository) Create(teacher *db.Teacher) error {
// 	if err := r.collection.Create(teacher); err != nil {
// 		return errors.New("failed to create teacher")
// 	}
// 	return nil
// }

// // FindByID retrieves a teacher by their teacherID
// func (r *MongoTeacherRepository) FindByID(teacherID int) (*db.Teacher, error) {
// 	teacher := &db.Teacher{}
// 	err := mgm.Coll(teacher).First(bson.M{"teacher_id": teacherID}, teacher)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return teacher, nil
// }

// // Update modifies an existing teacher's data
// func (r *MongoTeacherRepository) Update(teacher *db.Teacher) error {
// 	if err := r.collection.Update(teacher); err != nil {
// 		return errors.New("failed to update teacher")
// 	}
// 	return nil
// }

// // Delete removes a teacher from the database
// func (r *MongoTeacherRepository) Delete(teacherID int) error {
// 	teacher, err := r.FindByID(teacherID)
// 	if err != nil {
// 		return err
// 	}
// 	if err := r.collection.Delete(teacher); err != nil {
// 		return errors.New("failed to delete teacher")
// 	}
// 	return nil
// }

// func (r *MongoTeacherRepository) FindAll() ([]*db.Teacher, error) {
// 	var teachers []*db.Teacher
// 	if err := r.collection.SimpleFind(&teachers, bson.M{}); err != nil {
// 		return nil, errors.New("failed to retrieve teachers")
// 	}
// 	return teachers, nil
// }

// func (r *MongoTeacherRepository) FindByName(userName string) (*db.Teacher, error) {
// 	teacher := &db.Teacher{}
// 	err := mgm.Coll(teacher).First(bson.M{"name": userName}, teacher)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return teacher, nil
// }
