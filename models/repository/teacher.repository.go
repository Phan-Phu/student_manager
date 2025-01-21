package repository

import (
	"context"
	"errors"
	"student_manager/models"
	"student_manager/models/db"
	"student_manager/models/indexing"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeacherRepositoryInterface interface {
	Create(Teacher *db.Teacher) error
	FindByID(TeacherID string) (*db.Teacher, error)
	Update(Teacher *db.Teacher) error
	Delete(TeacherID string) error
	FindAll() ([]*db.Teacher, error)
	FindByFindOptions(findOptions *options.FindOptions) ([]db.Teacher, error)
	Count() (int64, error)
	FindByName(name string) (*db.Teacher, error)
	UpdateClass(class db.Class) (*db.Teacher, error)
	CreateMany(Teachers []db.Teacher) error
	DeleteMany(Teachers []db.Teacher) error
	GetTeacherDetails(TeacherIDString string) ([]*db.Class, error)
	GetTeachersDetails(TeacherIdsString []string) (*[]schemas.TeacherDetails, error)
}

type MongoTeacherRepository struct {
	TeacherCollection *mgm.Collection
}

func NewMongoTeacherRepository() *MongoTeacherRepository {
	indexing.NewIndexingTeacher()
	return &MongoTeacherRepository{
		TeacherCollection: mgm.Coll(&db.Teacher{}),
	}
}

func (r *MongoTeacherRepository) Create(Teacher *db.Teacher) error {
	err := r.TeacherCollection.Create(Teacher)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateTeacher))
	}

	return nil
}

func (r *MongoTeacherRepository) CreateMany(Teachers []db.Teacher) error {

	docs := make([]interface{}, len(Teachers))
	for i, Teacher := range Teachers {
		docs[i] = Teacher
	}

	_, err := r.TeacherCollection.InsertMany(mgm.Ctx(), docs)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateTeacher))
	}
	return nil
}

func (r *MongoTeacherRepository) DeleteMany(Teachers []interface{}) error {
	_, err := r.TeacherCollection.DeleteMany(mgm.Ctx(), Teachers)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateTeacher))
	}
	return nil
}

// FindByID retrieves a Teacher by their TeacherID
func (r *MongoTeacherRepository) FindByID(TeacherID string) (*db.Teacher, error) {
	Teacher := &db.Teacher{}
	objID, err := primitive.ObjectIDFromHex(TeacherID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}
	err = mgm.Coll(Teacher).First(bson.M{"_id": objID}, Teacher)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeFailRetrieveTeacher))
	}
	return Teacher, nil
}

func (r *MongoTeacherRepository) Update(Teacher *db.Teacher) error {
	if err := r.TeacherCollection.Update(Teacher); err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateTeacher))
	}
	return nil
}

func (r *MongoTeacherRepository) UpdateClass(class *db.Class) error {
	// Teachers, err := r.FindAll()
	// if err != nil {
	// 	return err
	// }

	// updated := false
	// for i := 0; i < len(Teachers); i++ {
	// 	if Teachers[i].ClassIds == class.ID {
	// 		Teachers[i].ClassIds = class.ID
	// 		updated = true
	// 	}
	// }

	// if !updated {
	// 	return errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	// }

	return nil
}

// Delete removes a Teacher from the database
func (r *MongoTeacherRepository) Delete(TeacherID string) error {
	Teacher, err := r.FindByID(TeacherID)
	if err != nil {
		return err
	}
	if err := r.TeacherCollection.Delete(Teacher); err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailDeleteTeacher))
	}
	return nil
}

func (r *MongoTeacherRepository) FindAll() ([]*db.Teacher, error) {
	var Teachers []*db.Teacher
	if err := r.TeacherCollection.SimpleFind(&Teachers, bson.M{}); err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeFailRetrieveTeacher))
	}
	return Teachers, nil
}

func (r *MongoTeacherRepository) FindByFindOptions(findOptions *options.FindOptions) ([]db.Teacher, error) {
	var Teachers []db.Teacher
	err := mgm.Coll(&db.Teacher{}).SimpleFind(&Teachers, bson.M{}, findOptions)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeTeacherIsNotFound))
	}
	return Teachers, nil
}

func (r *MongoTeacherRepository) Count() (int64, error) {
	total, err := mgm.Coll(&db.Teacher{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err != nil {
		return 0, errors.New(models.GetErrorMessage(models.ErrorCodeTeacherIsEmpty))
	}
	return total, nil
}

func (r *MongoTeacherRepository) FindByName(name string) (*db.Teacher, error) {
	Teacher := &db.Teacher{}
	err := mgm.Coll(Teacher).First(bson.M{"name": name}, Teacher)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeTeacherIsNotFound))
	}
	return Teacher, nil
}

func (r *MongoTeacherRepository) GetTeacherDetail(TeacherID primitive.ObjectID) (*schemas.TeacherDetails, error) {
	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: TeacherID},
			}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "classes"},
				{Key: "localField", Value: "class_ids"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "class"},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "birth_day", Value: 1},
				{Key: "class", Value: 1},
				{Key: "name", Value: 1},
			}},
		},
	}

	cursor, err := r.TeacherCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var TeacherResult []struct {
		Class    []*schemas.ClassResponse `bson:"class"`
		BirthDay string                   `bson:"birth_day"`
		Name     string                   `bson:"name"`
	}

	// result := []bson.M{}

	if err := cursor.All(context.Background(), &TeacherResult); err != nil {
		return nil, err
	}

	if len(TeacherResult) == 0 {
		return nil, nil
	}

	response := &schemas.TeacherDetails{
		TeacherID: TeacherID,
		Name:      TeacherResult[0].Name,
		Classes:   TeacherResult[0].Class,
		BirthDay:  TeacherResult[0].BirthDay,
	}

	return response, nil
}

func (r *MongoTeacherRepository) GetTeachersDetails(TeacherIds []primitive.ObjectID) ([]schemas.TeacherDetails, error) {

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "$in",
								Value: TeacherIds,
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "classes"},
					{Key: "localField", Value: "class_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "class"},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}, {Key: "birth_day", Value: 1}, {Key: "class", Value: 1}},
			},
		},
	}

	cursor, err := r.TeacherCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var TeacherResults []struct {
		Classes  []*schemas.ClassResponse `bson:"classes"`
		BirthDay string                   `bson:"birth_day"`
		Name     string                   `bson:"name"`
		ID       primitive.ObjectID       `bson:"_id"`
	}

	// result := []bson.M{}

	if err := cursor.All(context.Background(), &TeacherResults); err != nil {
		return nil, err
	}

	if len(TeacherResults) == 0 {
		return nil, nil
	}

	responses := []schemas.TeacherDetails{}

	for i := 0; i < len(TeacherResults); i++ {
		response := &schemas.TeacherDetails{
			TeacherID: TeacherResults[i].ID,
			Name:      TeacherResults[i].Name,
			Classes:   TeacherResults[i].Classes,
			BirthDay:  TeacherResults[i].BirthDay,
		}
		responses = append(responses, *response)
	}

	return responses, nil
}
