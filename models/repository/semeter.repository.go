package repository

// import (
// 	"student_manager/models/db"

// 	"github.com/kamva/mgm/v3"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type SemesterRepositoryInterface interface {
// 	Create(semester *db.Semester) error
// 	Update(semester *db.Semester) error
// 	FindById(studentID string) (*db.Semester, error)
// 	FindAll() ([]*db.Semester, error)
// }

// type SemesterRepository struct {
// 	collection *mgm.Collection
// }

// func NewSemesterRepository() *SemesterRepository {
// 	return &SemesterRepository{
// 		collection: mgm.Coll(&db.Semester{}),
// 	}
// }

// func (r *SemesterRepository) Create(semester *db.Semester) error {
// 	return r.collection.Create(semester)
// }

// func (r *SemesterRepository) Update(semester *db.Semester) error {
// 	return r.collection.Update(semester)
// }

// func (r *SemesterRepository) FindById(id primitive.ObjectID) (*db.Semester, error) {
// 	semester := &db.Semester{}
// 	err := r.collection.FindByID(id, semester)
// 	return semester, err
// }

// func (r *SemesterRepository) FindAll() ([]db.Semester, error) {
// 	var semesters []db.Semester
// 	err := r.collection.SimpleFind(&semesters, bson.M{})
// 	return semesters, err
// }
