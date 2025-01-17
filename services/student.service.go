package services

import (
	"errors"
	"fmt"
	"math"
	"studenent_manager/models"
	db "studenent_manager/models/db"
	"studenent_manager/models/indexing"
	"studenent_manager/models/repository"
	"studenent_manager/models/schemas"
	"sync"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	indexModels := []mongo.IndexModel{}
	// indexing.DeleteIndexing("age_1")
	indexing.DeleteAllIndexing()

	count, err := StudentService.Count()
	if err == nil && count > 0 {
		indexModels = append(indexModels, indexing.NewIndexingByStudentName())
		indexModels = append(indexModels, indexing.NewIndexingByStudentAge())
		indexModels = append(indexModels, indexing.NewPartialIndexScore(50)) // target score: 50
		indexModels = append(indexModels, indexing.NewIndexingByStudentAgeAndName())
		indexModels = append(indexModels, indexing.NewWildCardIndex())
		indexing.AddIndexModels(indexModels)
		indexing.PrintAllIndexing()

	}
}

func (repo *StudentRepoService) CreateStudent(data schemas.RequestStudent) (*schemas.StudentResponse, error) {
	initScore := 0

	classId, err := primitive.ObjectIDFromHex(data.ClassId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeInputIsWrong))
	}

	class, err := repository.NewMongoClassRepository().FindByID(data.ClassId)
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

func (repo *StudentRepoService) InsertManyStudentMultiThreads() error {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("InsertManyStudent took %s\n", elapsed)
	}()

	const totalStudents = 1000000
	const numThreads = 16
	batchSize := totalStudents / numThreads

	var wg sync.WaitGroup
	errChan := make(chan error, numThreads)
	threadTimes := make([]time.Duration, numThreads)

	// Worker function to process and insert students in batches
	worker := func(threadID, start, end int) {
		defer wg.Done()

		threadStart := time.Now()

		students := make([]db.Student, end-start)
		for i := start; i < end; i++ {
			// name := fmt.Sprintf("Name%d", i+1)
			// className := fmt.Sprintf("class %d", i+1)
			// student := db.NewStudent(name, className, "", 10, 0)
			// students[i-start] = *student
		}

		if err := repo.CreateMany(students); err != nil {
			errChan <- err
		}

		threadTimes[threadID] = time.Since(threadStart)
	}

	// Start worker threads
	for t := 0; t < numThreads; t++ {
		start := t * batchSize
		end := start + batchSize
		if t == numThreads-1 {
			end = totalStudents // Handle any remainder
		}

		wg.Add(1)
		go worker(t, start, end)
	}

	// Wait for all worker threads to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return errors.New(models.GetErrorMessage(models.ErrorCodeFailCreateStudent) + err.Error())
		}
	}

	// Print thread execution times
	for threadID, threadTime := range threadTimes {
		fmt.Printf("Thread %d took %s\n", threadID, threadTime)
	}

	return nil
}

func (repo *StudentRepoService) InsertManyStudentOneThread() error {
	start := time.Now()

	students := []db.Student{}

	for i := 0; i < 1000; i++ {
		// name := fmt.Sprintf("Name%d", i+1)
		// className := fmt.Sprintf("class %d", i+1)
		// student := db.NewStudent(name, className, "", 10, 0)
		// students = append(students, *student)
	}

	err := repo.CreateMany(students)
	if err != nil {
		return fmt.Errorf(models.GetErrorMessage(models.ErrorCodeFailCreateStudent), err)
	}

	estTime := time.Since(start)
	fmt.Printf("InsertManyStudentOneThreads took %s\n", estTime)

	return nil
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

	students, _ := repo.FindByFindOptions(findOptions)
	total, _ := repo.Count()

	paginationResponse := &schemas.PaginationResponse{
		Page:      data.Page,
		Limit:     data.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(data.Limit))),
	}

	studentIDs := funk.Map(students, func(s db.Student) string {
		return s.ID.Hex()
	}).([]string)

	studentDetails, err := repo.GetStudentsDetails(studentIDs)
	if err != nil {
		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	}

	studentResponses := funk.Map(studentDetails, func(s schemas.StudentDetails) schemas.StudentResponse {
		return *schemas.MapStudentDetailToStudentResponse(&s)
	}).([]schemas.StudentResponse)

	// studentResponses := make([]schemas.StudentResponse, len(students))
	// for i, student := range students {
	// 	details, err := repo.GetStudentDetails(student.ID.Hex())
	// 	if err != nil {
	// 		return nil, nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	// 	}
	// 	studentResponses[i] = *schemas.MapStudentDetailToStudentResponse(details)
	// }

	studentsResponse := &schemas.PaginationData[schemas.StudentResponse]{
		Data: studentResponses,
	}

	return paginationResponse, studentsResponse, nil
}

func (repo *StudentRepoService) GetStudent(studentId string) (*schemas.StudentResponse, error) {
	// student, err := repo.FindByID(studentName)
	// if err != nil {
	// 	return nil, errors.New(models.GetErrorMessage(models.ErrorCodeFailRetrieveStudent))
	// }

	details, err := repo.GetStudentDetails(studentId)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeGetStudentDetails))
	}

	response := schemas.MapStudentDetailToStudentResponse(details)

	return response, nil
}

func (repo *StudentRepoService) UpdateStudent(data schemas.UpdateStudent) (*schemas.ResponseUpdateStudent, error) {
	student, err := repo.FindByID(data.StudentID)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeMaxStudentInClass))
	}

	student.Name = data.Name
	// student.ClassName = data.ClassName
	student.BirthDay = data.BirthDay
	student.Age = data.Age

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeCanNotUpdateStudent))
	}

	details, err := repo.GetStudentDetails(student.ID.Hex())
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
