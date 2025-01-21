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

var TeacherService *TeacherRepoService

type TeacherRepoService struct {
	*repository.MongoTeacherRepository
}

func InitializeTeacherRepository() {
	TeacherService = &TeacherRepoService{
		MongoTeacherRepository: repository.NewMongoTeacherRepository(),
	}
}

func (repo *TeacherRepoService) CreateTeacher(data schemas.RequestTeacher) (*schemas.TeacherResponse, error) {

	classIds, err := ConvertStringsToObjectIDs(data.ClassIds)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	for i := 0; i < len(classIds); i++ {
		if _, err := ClassService.FindByID(classIds[i]); err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeNotFoundClass))
		}
	}

	teacher := db.NewTeacher(data.Name, classIds, data.BirthDay, data.Username, data.Password)
	err = repo.Create(teacher)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotCreateTeacher))
	}

	teacherDetails, err := repo.GetTeacherDetail(teacher.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetTeacherDetails))
	}

	response := &schemas.TeacherResponse{
		TeacherID: teacherDetails.TeacherID,
		Name:      teacherDetails.Name,
		Classes:   teacherDetails.Classes,
		BirthDay:  teacher.BirthDay,
	}

	return response, nil
}

func (repo *TeacherRepoService) GetTeachers() ([]*db.Teacher, error) {
	Teachers, err := repo.FindAll()

	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxTeacherInClass))
	}
	return Teachers, nil
}

func (repo *TeacherRepoService) GetTeachersWithPagination(data *schemas.PaginationRequest) (*schemas.PaginationResponse, *schemas.PaginationData[schemas.TeacherResponse], error) {
	findOptions := options.Find()
	findOptions.SetSkip((int64(data.Page) - 1) * data.Limit)
	findOptions.SetLimit(int64(data.Limit))

	Teachers, err := repo.FindByFindOptions(findOptions)
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	total, err := repo.Count()
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxTeacherInClass))
	}

	paginationResponse := &schemas.PaginationResponse{
		Page:      data.Page,
		Limit:     data.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(data.Limit))),
	}

	TeacherIDs := []primitive.ObjectID{}
	for i := 0; i < len(Teachers); i++ {
		TeacherIDs = append(TeacherIDs, Teachers[i].ID)
	}

	TeacherDetails, err := repo.GetTeachersDetails(TeacherIDs)
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetTeacherDetails))
	}

	TeacherResponses := []schemas.TeacherResponse{}
	for i := 0; i < len(TeacherDetails); i++ {
		TeacherResponses = append(TeacherResponses, *schemas.MapTeacherDetailToTeacherResponse(&TeacherDetails[i]))
	}

	TeachersResponse := &schemas.PaginationData[schemas.TeacherResponse]{
		Data: TeacherResponses,
	}

	return paginationResponse, TeachersResponse, nil
}

func (repo *TeacherRepoService) GetTeacher(TeacherIDString string) (*schemas.TeacherResponse, error) {
	TeacherId, err := primitive.ObjectIDFromHex(TeacherIDString)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	details, err := repo.GetTeacherDetail(TeacherId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetTeacherDetails))
	}

	response := schemas.MapTeacherDetailToTeacherResponse(details)

	return response, nil
}

func (repo *TeacherRepoService) UpdateInfoTeacher(data schemas.UpdateInfoTeacherRequest) (*schemas.ResponseUpdateTeacher, error) {
	teacher, err := repo.FindByID(data.TeacherID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeNotFoundTeacher))
	}

	teacher.Name = data.Name
	teacher.BirthDay = data.BirthDay

	err = repo.Update(teacher)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateTeacher))
	}

	details, err := repo.GetTeacherDetail(teacher.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetTeacherDetails))
	}

	response := &schemas.ResponseUpdateTeacher{
		TeacherID: teacher.ID,
		Name:      teacher.Name,
		Classes:   details.Classes,
		BirthDay:  teacher.BirthDay,
	}

	return response, nil
}

func (repo *TeacherRepoService) UpdateClassTeacher(data schemas.UpdateTeacherClassRequest) (*schemas.ResponseUpdateTeacher, error) {
	teacher, err := repo.FindByID(data.TeacherID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	if len(data.ClassesUpdate) == 0 {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeUpdateClassIsEmpty))
	}

	newClassIds := make([]primitive.ObjectID, len(data.ClassesUpdate))
	for i, classUpdate := range data.ClassesUpdate {
		oldClassId, err := ConvertStringToObjectID(classUpdate.OldClassID)
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrongTeacher))
		}

		if !contains(teacher.ClassIds, *oldClassId) {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeTeacherNotAssignedToClass))
		}

		_, err = ClassService.FindByID(*oldClassId)
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindClass))
		}

		newClassId, err := ConvertStringToObjectID(classUpdate.NewClassID)
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrongTeacher))
		}
		newClassIds[i] = *newClassId
	}

	teacher.ClassIds = newClassIds
	err = repo.Update(teacher)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateTeacher))
	}

	details, err := repo.GetTeacherDetail(teacher.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetTeacherDetails))
	}

	response := &schemas.ResponseUpdateTeacher{
		TeacherID: teacher.ID,
		Name:      teacher.Name,
		Classes:   details.Classes,
		BirthDay:  teacher.BirthDay,
	}

	return response, nil
}

func contains(ids []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func (repo *TeacherRepoService) DeleteTeacher(data schemas.DeleteTeacher) error {
	Teacher, err := repo.FindByID(data.TeacherID)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeMaxTeacherInClass))
	}

	err = mgm.Coll(Teacher).Delete(Teacher)
	if err != nil {
		return errors.New(models.GetErrorMessage(models.ErrorCodeFailDeleteTeacher))
	}

	return nil
}
