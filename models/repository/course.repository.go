package repository

import (
	"context"
	"student_manager/models/db"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepositoryInterface interface {
	Create(course *db.Course) error
	Update(course *db.Course) error
	FindById(id primitive.ObjectID) (*db.Course, error)
	FindAll() ([]db.Course, error)
}

type CourseRepository struct {
	collection *mgm.Collection
}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{
		collection: mgm.Coll(&db.Course{}),
	}
}

func (r *CourseRepository) Create(course *db.Course) error {
	return r.collection.Create(course)
}

func (r *CourseRepository) Update(course *db.Course) error {
	return r.collection.Update(course)
}

func (r *CourseRepository) FindById(id primitive.ObjectID) (*db.Course, error) {
	course := &db.Course{}
	err := r.collection.First(bson.M{"_id": id}, course)
	return course, err
}

func (r *CourseRepository) FindAll() ([]db.Course, error) {
	courses := []db.Course{}
	err := r.collection.SimpleFind(&courses, bson.M{})
	return courses, err
}

func (r *CourseRepository) GetCourseDetails(enrollmentId primitive.ObjectID) (*schemas.CourseDetail, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"_id": enrollmentId,
			}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "students",
				"localField":   "student_ids",
				"foreignField": "_id",
				"as":           "student_details",
			}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "teachers",
				"localField":   "teacher_ids",
				"foreignField": "_id",
				"as":           "teacher_details",
			}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "coursees",
				"localField":   "course_id",
				"foreignField": "_id",
				"as":           "course_details",
			}},
		},
		{
			{Key: "$unwind", Value: "$course_details"},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "subjects",
				"localField":   "subject_id",
				"foreignField": "_id",
				"as":           "subject_details",
			}},
		},
		{
			{Key: "$unwind", Value: "$subject_details"},
		},
		{
			{Key: "$project", Value: bson.M{
				"status":   1,
				"students": "$student_details",
				"teachers": "$teacher_details",
				"course":   "$course_details",
				"subject":  "$subject_details",
			}},
		},
	}

	result := []schemas.CourseDetail{}
	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return &result[0], nil
}
