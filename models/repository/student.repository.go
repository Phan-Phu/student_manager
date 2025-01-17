package repository

import (
	"context"
	"errors"
	"studenent_manager/models"
	"studenent_manager/models/db"
	"studenent_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentRepositoryInterface interface {
	Create(student *db.Student) error
	FindByID(studentID string) (*db.Student, error)
	update(student *db.Student) error
	Delete(studentID string) error
	FindAll() ([]*db.Student, error)
	FindByFindOptions(findOptions *options.FindOptions) ([]db.Student, error)
	Count() (int64, error)
	FindByName(name string) (*db.Student, error)
	UpdateClass(class db.Class) (*db.Student, error)
	CreateMany(students []db.Student) error
	DeleteMany(students []db.Student) error
	UpdateWithClass() error
	GetStudentDetails(studentIDString string) ([]*db.Class, error)
	GetStudentsDetails(studentIdsString []string) (*[]schemas.StudentDetails, error)
}

type MongoStudentRepository struct {
	studentCollection *mgm.Collection
}

func NewMongoStudentRepository() *MongoStudentRepository {
	return &MongoStudentRepository{
		studentCollection: mgm.Coll(&db.Student{}),
	}
}

func (r *MongoStudentRepository) Create(student *db.Student) error {
	err := r.studentCollection.Create(student)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent))
	}

	return nil
}

func (r *MongoStudentRepository) CreateMany(students []db.Student) error {

	docs := make([]interface{}, len(students))
	for i, student := range students {
		docs[i] = student
	}

	_, err := r.studentCollection.InsertMany(mgm.Ctx(), docs)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent))
	}
	return nil
}

func (r *MongoStudentRepository) DeleteMany(students []interface{}) error {
	_, err := r.studentCollection.DeleteMany(mgm.Ctx(), students)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent))
	}
	return nil
}

// FindByID retrieves a student by their StudentID
func (r *MongoStudentRepository) FindByID(studentID string) (*db.Student, error) {
	student := &db.Student{}
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}
	err = mgm.Coll(student).First(bson.M{"_id": objID}, student)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeFailRetrieveStudent))
	}
	return student, nil
}

// Update modifies an existing student's data
func (r *MongoStudentRepository) Update(student *db.Student) error {
	if err := r.studentCollection.Update(student); err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent))
	}
	return nil
}

func (r *MongoStudentRepository) UpdateClass(class *db.Class) error {
	students, err := r.FindAll()
	if err != nil {
		return err
	}

	updated := false
	for i := 0; i < len(students); i++ {
		if students[i].ClassId == class.ID {
			students[i].ClassId = class.ID
			updated = true
		}
	}

	if !updated {
		return errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}

	return nil
}

// Delete removes a student from the database
func (r *MongoStudentRepository) Delete(studentID string) error {
	student, err := r.FindByID(studentID)
	if err != nil {
		return err
	}
	if err := r.studentCollection.Delete(student); err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailDeleteStudent))
	}
	return nil
}

func (r *MongoStudentRepository) FindAll() ([]*db.Student, error) {
	var students []*db.Student
	if err := r.studentCollection.SimpleFind(&students, bson.M{}); err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeFailRetrieveStudent))
	}
	return students, nil
}

func (r *MongoStudentRepository) FindByFindOptions(findOptions *options.FindOptions) ([]db.Student, error) {
	var students []db.Student
	err := mgm.Coll(&db.Student{}).SimpleFind(&students, bson.M{}, findOptions)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeStudentIsNotFound))
	}
	return students, nil
}

func (r *MongoStudentRepository) Count() (int64, error) {
	total, err := mgm.Coll(&db.Student{}).CountDocuments(mgm.Ctx(), bson.M{})
	if err != nil {
		return 0, errors.New(models.GetErrorMessage(models.ErrorCodeStudentIsEmpty))
	}
	return total, nil
}

func (r *MongoStudentRepository) FindByName(name string) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"name": name}, student)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeStudentIsNotFound))
	}
	return student, nil
}

func (r *MongoStudentRepository) GetStudentDetails(studentIDString string) (*schemas.StudentDetails, error) {
	studentId, err := primitive.ObjectIDFromHex(studentIDString)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: studentId},
			}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "classes"},
				{Key: "localField", Value: "class_id"},
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

	cursor, err := r.studentCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var studentResult []struct {
		Class    []*schemas.ClassResponse `bson:"class"`
		BirthDay string                   `bson:"birth_day"`
		Name     string                   `bson:"name"`
	}

	// result := []bson.M{}

	if err := cursor.All(context.Background(), &studentResult); err != nil {
		return nil, err
	}

	if len(studentResult) == 0 || len(studentResult[0].Class) == 0 {
		return nil, nil
	}

	response := &schemas.StudentDetails{
		StudentID: studentId,
		Name:      studentResult[0].Name,
		Class:     *studentResult[0].Class[0],
		BirthDay:  studentResult[0].BirthDay,
	}

	return response, nil
}

func (r *MongoStudentRepository) GetStudentsDetails(studentIdsString []string) ([]schemas.StudentDetails, error) {

	studentIds := []primitive.ObjectID{}
	for i := 0; i < len(studentIdsString); i++ {
		studentId, err := primitive.ObjectIDFromHex(studentIdsString[i])
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
		}
		studentIds = append(studentIds, studentId)
	}

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "$in",
								Value: studentIds,
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

	cursor, err := r.studentCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var studentResults []struct {
		Class    []*schemas.ClassResponse `bson:"class"`
		BirthDay string                   `bson:"birth_day"`
		Name     string                   `bson:"name"`
		ID       primitive.ObjectID       `bson:"_id"`
	}

	// result := []bson.M{}

	if err := cursor.All(context.Background(), &studentResults); err != nil {
		return nil, err
	}

	if len(studentResults) == 0 {
		return nil, nil
	}

	responses := []schemas.StudentDetails{}

	for i := 0; i < len(studentResults); i++ {
		response := &schemas.StudentDetails{
			StudentID: studentResults[i].ID,
			Name:      studentResults[i].Name,
			Class:     *studentResults[i].Class[0],
			BirthDay:  studentResults[i].BirthDay,
		}
		responses = append(responses, *response)
	}

	return responses, nil
}
