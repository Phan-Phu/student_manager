package services

import (
	"errors"
	"student_manager/models"
	"student_manager/models/db"
	"student_manager/models/repository"
	"student_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ClassService *ClassRepoService

type ClassRepoService struct {
	*repository.ClassRepository
}

func InitializeClassRepository() {
	ClassService = &ClassRepoService{
		ClassRepository: repository.NewMongoClassRepository(),
	}
}

func (repo *ClassRepoService) CreateClass(data schemas.RequestClass) (*schemas.ClassResponse, error) {

	class := db.NewClass(data.Name, data.MaxStudent)
	err := repo.Create(class)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotCreateClass))
	}

	classResponse := &schemas.ClassResponse{
		Name:       class.Name,
		MaxStudent: class.MaxStudent,
	}

	return classResponse, nil
}

func (repo *ClassRepoService) GetClasses() ([]*db.Class, error) {
	classes, err := repo.FindAll()

	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}
	return classes, nil
}

func (repo *ClassRepoService) GetClass(classIDString string) (*schemas.ClassResponse, error) {
	classId, err := primitive.ObjectIDFromHex(classIDString)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	details, err := repo.FindByID(classId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}

	response := &schemas.ClassResponse{
		ClassID:    details.ID,
		CreateDate: details.CreateDate,
		Name:       details.Name,
		MaxStudent: details.MaxStudent,
	}

	return response, nil
}

func (repo *ClassRepoService) UpdateClass(data schemas.UpdateClassRequest) (*schemas.ResponseClassUpdate, error) {
	classId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	class, err := repo.FindByID(classId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}

	class.Name = data.Name
	class.MaxStudent = data.MaxStudent

	err = ClassService.Update(class)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateClass))
	}

	response := &schemas.ResponseClassUpdate{
		ID:         class.ID.Hex(),
		Name:       class.Name,
		MaxStudent: class.MaxStudent,
	}

	return response, nil
}

func (repo *ClassRepoService) DeleteClass(data schemas.DeleteClassRequest) error {
	classId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	Class, err := repo.FindByID(classId)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
	}

	err = mgm.Coll(Class).Delete(Class)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailDeleteClass))
	}

	return nil
}
