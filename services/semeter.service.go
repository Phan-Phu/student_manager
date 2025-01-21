package services

// import (
// 	"errors"
// 	"student_manager/models"
// 	"student_manager/models/db"
// 	"student_manager/models/repository"
// 	"student_manager/models/schemas"
// )

// var SemesterService *SemesterRepoService

// type SemesterRepoService struct {
// 	repo *repository.SemesterRepository
// }

// func InitializeSemesterService() {
// 	SemesterService = &SemesterRepoService{
// 		repo: repository.NewSemesterRepository(),
// 	}
// }

// func (s *SemesterRepoService) CreateSemester(req schemas.CreateSemesterRequest) (*schemas.SemesterResponse, error) {
// 	semester := &db.Semester{
// 		Name:       req.Name,
// 		CreateDate: req.CreateDate,
// 		EndDate:    req.EndDate,
// 		Classes:    []db.Course{},
// 	}

// 	if err := s.repo.Create(semester); err != nil {
// 		return nil, err
// 	}

// 	return &schemas.SemesterResponse{
// 		ID:         semester.ID.Hex(),
// 		Name:       semester.Name,
// 		CreateDate: semester.CreateDate,
// 		EndDate:    semester.EndDate,
// 		Classes:    []schemas.CourseResponse{},
// 	}, nil
// }

// func (s *SemesterRepoService) UpdateSemester(req schemas.UpdateSemesterRequest) (*schemas.SemesterResponse, error) {
// 	semesterId, err := ConvertStringToObjectID(req.ID)
// 	if err != nil {
// 		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
// 	}

// 	existingSemester, err := s.repo.FindById(*semesterId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	existingSemester.Name = req.Name
// 	existingSemester.CreateDate = req.CreateDate
// 	existingSemester.EndDate = req.EndDate

// 	if err := s.repo.Update(existingSemester); err != nil {
// 		return nil, err
// 	}

// 	return &schemas.SemesterResponse{
// 		ID:         existingSemester.ID.Hex(),
// 		Name:       existingSemester.Name,
// 		CreateDate: existingSemester.CreateDate,
// 		EndDate:    existingSemester.EndDate,
// 		Classes:    []schemas.CourseResponse{},
// 	}, nil
// }
