package services

import (
	"errors"
	"math"
	"student_manager/models"
	db "student_manager/models/db"
	"student_manager/models/repository"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var StudentService *StudentRepoService

type StudentRepoService struct {
	*repository.MongoStudentRepository
}

func InitializeStudentRepository() {
	StudentService = &StudentRepoService{
		MongoStudentRepository: repository.NewMongoStudentRepository(),
	}
}

func (repo *StudentRepoService) CreateStudent(data schemas.RequestStudent) (*schemas.StudentResponse, error) {
	initScore := 0

	classId, err := primitive.ObjectIDFromHex(data.ClassId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	class, err := ClassService.FindByID(classId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	if class.MaxStudent > db.MaxStudentInClass {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxStudentInClass))
	}

	student := db.NewStudent(data.Name, classId, data.BirthDay, data.Age, initScore)
	err = repo.Create(student)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotCreateStudent))
	}

	classResponse := &schemas.ClassResponse{
		Name:       class.Name,
		MaxStudent: class.MaxStudent,
	}

	response := schemas.MapStudentToStudentResponse(student, classResponse)

	return response, nil
}

func (repo *StudentRepoService) GetStudents() ([]*db.Student, error) {
	students, err := repo.FindAll()

	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxStudentInClass))
	}
	return students, nil
}

func (repo *StudentRepoService) GetStudentsWithPagination(data *schemas.PaginationRequest) (*schemas.PaginationResponse, *schemas.PaginationData[schemas.StudentResponse], error) {
	findOptions := options.Find()
	findOptions.SetSkip((int64(data.Page) - 1) * data.Limit)
	findOptions.SetLimit(int64(data.Limit))

	students, err := repo.FindByFindOptions(findOptions)
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	total, err := repo.Count()
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxStudentInClass))
	}

	paginationResponse := &schemas.PaginationResponse{
		Page:      data.Page,
		Limit:     data.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(data.Limit))),
	}

	studentIDs := []primitive.ObjectID{}
	for i := 0; i < len(students); i++ {
		studentIDs = append(studentIDs, students[i].ID)
	}

	studentDetails, err := repo.GetStudentsDetails(studentIDs)
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	}

	studentResponses := []schemas.StudentResponse{}
	for i := 0; i < len(studentDetails); i++ {
		studentResponses = append(studentResponses, *schemas.MapStudentDetailToStudentResponse(&studentDetails[i]))
	}

	studentsResponse := &schemas.PaginationData[schemas.StudentResponse]{
		Data: studentResponses,
	}

	return paginationResponse, studentsResponse, nil
}

func (repo *StudentRepoService) GetStudent(studentIDString string) (*schemas.StudentResponse, error) {
	studentId, err := primitive.ObjectIDFromHex(studentIDString)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	details, err := repo.GetStudentDetails(studentId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	}

	response := schemas.MapStudentDetailToStudentResponse(details)

	return response, nil
}

func (repo *StudentRepoService) UpdateStudent(studentId string, data schemas.UpdateStudentRequest) (*schemas.ResponseUpdateStudent, error) {
	student, err := repo.FindByID(studentId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	student.Name = data.Name
	// student.ClassName = data.ClassName
	student.BirthDay = data.BirthDay
	student.Age = data.Age

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateStudent))
	}

	details, err := repo.GetStudentDetails(student.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	}

	response := &schemas.ResponseUpdateStudent{
		StudentID: student.ID,
		Name:      student.Name,
		Class:     details.Class,
		BirthDay:  student.BirthDay,
	}

	return response, nil
}

func (repo *StudentRepoService) UpdateScore(data schemas.UpdateScoreStudent) (*schemas.ResponseUpdateScoreStudent, error) {
	student, err := repo.FindByID(data.StudentID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeStudentIsNotFound))
	}

	student.Score = data.Score

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New((models.GetErrorMessage(models.ErrorCodeCanNotUpdateStudent)))
	}

	response := &schemas.ResponseUpdateScoreStudent{
		StudentID:   student.ID,
		StudentName: student.Name,
		Score:       student.Score,
	}

	return response, nil
}

func (repo *StudentRepoService) DeleteStudent(data schemas.DeleteStudent) error {
	student, err := repo.FindByID(data.StudentID)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeMaxStudentInClass))
	}

	err = mgm.Coll(student).Delete(student)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailDeleteStudent))
	}

	return nil
}
