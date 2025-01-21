package repository

import (
	"context"
	"errors"
	"student_manager/models"
	"student_manager/models/db"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreRepositoryInterface interface {
	Create(score *db.Score) error
	CreateMany(scores []*db.Score) ([]primitive.ObjectID, error)
	Update(score *db.Score) error
	FindById(scoreId primitive.ObjectID) (*db.Score, error)
	GetScoreResponse(scoreId primitive.ObjectID) (*schemas.ScoreDetails, error)
}

type ScoreRepository struct {
	collection *mgm.Collection
}

func NewScoreRepository() *ScoreRepository {
	return &ScoreRepository{
		collection: mgm.Coll(&db.Score{}),
	}
}

func (r *ScoreRepository) Create(score *db.Score) error {
	return r.collection.Create(score)
}

func (r *ScoreRepository) CreateMany(scores []*db.Score) ([]primitive.ObjectID, error) {
	var interfaceSlice []interface{}
	for _, scoreItem := range scores {
		interfaceSlice = append(interfaceSlice, scoreItem)
	}

	result, err := r.collection.InsertMany(context.Background(), interfaceSlice)
	if err != nil {
		return nil, err
	}
	objectIDs := []primitive.ObjectID{}

	for _, id := range result.InsertedIDs {
		if objectID, ok := id.(primitive.ObjectID); ok {
			objectIDs = append(objectIDs, objectID)
		} else {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
		}
	}

	return objectIDs, nil
}

func (r *ScoreRepository) Update(score *db.Score) error {
	return r.collection.Update(score)
}

func (r *ScoreRepository) FindById(scoreId primitive.ObjectID) (*db.Score, error) {
	score := &db.Score{}
	err := r.collection.First(bson.M{"_id": scoreId}, score)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindScoreInCourse))
	}
	return score, nil
}

func (r *ScoreRepository) GetScoreResponse(scoreId primitive.ObjectID) (*schemas.ScoreDetails, error) {
	pipeline := bson.A{
		bson.D{{
			Key:   "$match",
			Value: bson.D{{Key: "_id", Value: scoreId}},
		}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "subjects"},
					{Key: "localField", Value: "subject_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "subject"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "students"},
					{Key: "localField", Value: "student_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: `student`},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "classes"},
					{Key: `localField`, Value: "student.class_id"},
					{Key: `foreignField`, Value: `_id`},
					{Key: "as", Value: "class"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$class"},
					{Key: `preserveNullAndEmptyArrays`, Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$student"},
					{Key: `preserveNullAndEmptyArrays`, Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$subject"},
					{Key: `preserveNullAndEmptyArrays`, Value: true},
				},
			},
		},
	}

	detail := []schemas.ScoreDetails{}
	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &detail); err != nil {
		return nil, err
	}

	return &detail[0], nil
}
