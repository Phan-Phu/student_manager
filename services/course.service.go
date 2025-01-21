package services

import (
	"errors"
	"student_manager/models"
	"student_manager/models/db"
	"student_manager/models/repository"
	"student_manager/models/schemas"
)

var CourseService *CourseRepoService

type CourseRepoService struct {
	courseRepo *repository.CourseRepository
	scoreRepo  *repository.ScoreRepository
}

func InitializeCourseService() {
	CourseService = &CourseRepoService{
		courseRepo: repository.NewCourseRepository(),
		scoreRepo:  repository.NewScoreRepository(),
	}
}

func (s *CourseRepoService) CreateCourse(req schemas.CreateCourseRequest) (*schemas.CourseResponse, error) {
	studentIds, err := ConvertStringsToObjectIDs(req.StudentIds)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	teacherIds, err := ConvertStringsToObjectIDs(req.TeacherIds)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	classId, err := ConvertStringToObjectID(req.ClassId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	subjectId, err := ConvertStringToObjectID(req.SubjectId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	defaultScores := make(map[db.ScoreType]int)
	for _, exam := range req.Examinations {
		category := db.ScoreType(exam)
		defaultScores[category] = 0
	}

	scores := make([]*db.Score, len(studentIds))

	for i, studentId := range studentIds {
		scores[i] = db.NewScore(studentId, *subjectId, defaultScores)
	}

	scoresIds, err := s.scoreRepo.CreateMany(scores)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCreateManyScore))
	}

	course := db.NewCourse(req.Name, studentIds, teacherIds, *classId, *subjectId, req.StartDate, req.EndDate, scoresIds)

	if err := s.courseRepo.Create(course); err != nil {
		return nil, err
	}

	details, err := s.courseRepo.GetCourseDetails(course.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetCourseDetails))
	}

	// Tạo response
	response := &schemas.CourseResponse{
		ID:        course.ID.Hex(),
		Students:  details.Students,
		Teachers:  details.Teachers,
		Classes:   details.Class,
		Subject:   details.Subject,
		StartDate: course.StartDate,
		EndDate:   course.EndDate,
		ScoreIDs:  course.Scores,
	}

	return response, nil
}

func (s *CourseRepoService) UpdateCourse(id string, req schemas.UpdateCourseRequest) (*schemas.CourseResponse, error) {
	courseId, err := ConvertStringToObjectID(id)
	if err != nil {
		return nil, err
	}

	studentIds, err := ConvertStringsToObjectIDs(req.StudentIds)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	teacherIds, err := ConvertStringsToObjectIDs(req.TeacherIds)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	classId, err := ConvertStringToObjectID(req.ClassId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	subjectId, err := ConvertStringToObjectID(req.SubjectId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}

	course, err := s.courseRepo.FindById(*courseId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindCourse))
	}

	course.StudentIds = studentIds
	course.TeacherIds = teacherIds
	course.ClassId = *classId
	course.SubjectId = *subjectId
	course.StartDate = req.StartDate
	course.EndDate = req.EndDate

	err = s.courseRepo.Update(course)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateCourse))
	}

	details, err := s.courseRepo.GetCourseDetails(course.ID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetCourseDetails))
	}

	response := &schemas.CourseResponse{
		ID:        course.ID.Hex(),
		Name:      course.Name,
		Students:  details.Students,
		Teachers:  details.Teachers,
		Classes:   details.Class,
		Subject:   details.Subject,
		StartDate: course.StartDate,
		EndDate:   course.EndDate,
		ScoreIDs:  course.Scores,
	}

	return response, nil
}

func (s *CourseRepoService) UpdateScore(req schemas.UpdateScoreRequest) (*schemas.ScoreResponse, error) {
	id, err := ConvertStringToObjectID(req.ScoreId)
	if err != nil {
		return nil, err
	}

	score, err := s.scoreRepo.FindById(*id)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindScoreInCourse))
	}

	score.Score[db.ScoreType(req.Type)] = req.Score

	if err := s.scoreRepo.Update(score); err != nil {
		return nil, err
	}

	detail, err := s.scoreRepo.GetScoreResponse(*id)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetScoreDetails))
	}

	response := &schemas.ScoreResponse{
		ScoreId: detail.ScoreId,
		Student: detail.Student,
		Subject: detail.Subject,
		Score:   detail.Score,
	}
	response.Student.Class = detail.Class

	return response, nil
}

func (s *CourseRepoService) GetCourses() ([]*schemas.CourseResponse, error) {

	courses, err := s.courseRepo.FindAll()
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotFindScoreInCourse))
	}

	response := []*schemas.CourseResponse{}

	for i := 0; i < len(courses); i++ {
		detail, err := s.courseRepo.GetCourseDetails(courses[i].ID)
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetScoreDetails))
		}

		data := &schemas.CourseResponse{
			ID:        courses[i].ID.Hex(),
			Name:      courses[i].Name,
			Students:  detail.Students,
			Teachers:  detail.Teachers,
			Classes:   detail.Class,
			Subject:   detail.Subject,
			StartDate: courses[i].StartDate,
			EndDate:   courses[i].EndDate,
		}

		response = append(response, data)
	}

	return response, nil
}
