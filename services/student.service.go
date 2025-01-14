package services

import (
	"errors"
	"fmt"
	"math"
	db "studenent_manager/models/db"
	"studenent_manager/models/indexing"
	"studenent_manager/models/repository"
	"studenent_manager/models/schemas"
	"sync"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var StudentService *StudentRepoService

type StudentRepoService struct {
	Repo repository.StudentRepositoryInteface
}

func InitializeStudentRepository() {
	studentrepo := repository.NewMongoStudentRepository()
	StudentService = &StudentRepoService{
		Repo: studentrepo,
	}
	indexModels := []mongo.IndexModel{}
	// indexing.DeleteIndexing("age_1")
	indexing.DeleteAllIndexing()

	// not check if not first student --> error
	indexModels = append(indexModels, indexing.NewIndexingByStudentID())
	indexModels = append(indexModels, indexing.NewIndexingByStudentName())
	indexModels = append(indexModels, indexing.NewIndexingByStudentAge())
	indexModels = append(indexModels, indexing.NewPartialIndexScore(50)) // target score: 50
	indexModels = append(indexModels, indexing.NewIndexingByStudentAgeAndName())
	indexModels = append(indexModels, indexing.NewWildCardIndex())
	indexing.AddIndexModels(indexModels)
	indexing.PrintAllIndexing()
}

func (repo *StudentRepoService) CreateStudent(data schemas.RequestStudent) (*schemas.ResponseStudent, error) {
	initScore := 0

	student := db.NewStudent(data.Name, data.ClassID, data.BirthDay, data.Age, initScore)
	err := repo.Repo.Create(student)

	if err != nil {
		return nil, errors.New("cannot create new student")
	}

	response := &schemas.ResponseStudent{
		StudentID: student.ID,
		Name:      student.Name,
		ClassID:   student.ClassID,
		BirthDay:  student.BirthDay,
	}

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
			name := fmt.Sprintf("Name%d", i+1)
			student := db.NewStudent(name, i+1, "", 10, 0)
			students[i-start] = *student
		}

		docs := make([]interface{}, len(students))
		for i, student := range students {
			docs[i] = student
		}

		if err := repo.Repo.CreateMany(docs); err != nil {
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
			return errors.New("failed to create students: " + err.Error())
		}
	}

	// Print thread execution times
	for threadID, threadTime := range threadTimes {
		fmt.Printf("Thread %d took %s\n", threadID, threadTime)
	}

	return nil
}

func (repo *StudentRepoService) InsertManyStudentOneThread() error {
	start := time.Now() // Bắt đầu đo thời gian

	students := []db.Student{}

	// Tạo danh sách students
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("Name%d", i+1)
		student := db.NewStudent(name, i+1, "", 10, 0)
		students = append(students, *student)
	}

	// Chuyển đổi students thành danh sách interface{}
	docs := make([]interface{}, len(students))
	for i, student := range students {
		docs[i] = student
	}

	// Chèn dữ liệu vào database
	err := repo.Repo.CreateMany(docs)
	if err != nil {
		return fmt.Errorf("failed to insert students: %w", err)
	}

	// Đo thời gian sau khi hoàn thành toàn bộ công việc
	estTime := time.Since(start)
	fmt.Printf("InsertManyStudentOneThreads took %s\n", estTime)

	return nil
}

func (repo *StudentRepoService) GetStudents() ([]*db.Student, error) {
	students, err := repo.Repo.FindAll()

	if err != nil {
		return nil, errors.New("cannot get students")
	}
	return students, nil
}

func (repo *StudentRepoService) GetStudentsWithPagination(data *schemas.PaginationRequest) (*schemas.PaginationResponse, *schemas.PaginationData[db.Student], error) {
	findOptions := options.Find()
	findOptions.SetSkip((int64(data.Page) - 1) * data.Limit)
	findOptions.SetLimit(int64(data.Limit))

	students, _ := repo.Repo.FindByFindOptions(findOptions)
	total, _ := repo.Repo.Count()

	paginationResponse := &schemas.PaginationResponse{
		Page:      data.Page,
		Limit:     data.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(data.Limit))),
	}

	studentsRespone := &schemas.PaginationData[db.Student]{
		Data: students,
	}

	return paginationResponse, studentsRespone, nil
}

func (repo *StudentRepoService) GetStudent(studentName string) (*schemas.ResponseStudent, error) {
	start := time.Now()
	student, err := repo.Repo.FindByName(studentName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	elapsed := time.Since(start)
	fmt.Printf("Found student: %+v\n", student)
	fmt.Printf("Time with index: %v\n", elapsed)

	response := &schemas.ResponseStudent{
		StudentID: student.ID,
		Name:      student.Name,
		ClassID:   student.ClassID,
		Age:       student.Age,
		BirthDay:  student.BirthDay,
	}

	return response, nil
}

func (repo *StudentRepoService) UpdateStudent(data schemas.UpdateStudent) (*schemas.ResponseUpdateStudent, error) {
	student, err := repo.Repo.FindByID(data.StudentID)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	student.Name = data.Name
	student.ClassID = data.ClassID
	student.BirthDay = data.BirthDay
	student.Age = data.Age

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New("cannot update student")
	}

	response := &schemas.ResponseUpdateStudent{
		StudentID: student.ID,
		Name:      student.Name,
		ClassID:   student.ClassID,
		Age:       student.Age,
		BirthDay:  student.BirthDay,
	}

	return response, nil
}

func (repo *StudentRepoService) UpdateScore(data schemas.UpdateScoreStudent) (*schemas.ResponseUpdateScoreStudent, error) {
	student, err := repo.Repo.FindByID(data.StudentID)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	student.Score = data.Score

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New("cannot update student")
	}

	response := &schemas.ResponseUpdateScoreStudent{
		StudentID:   student.ID,
		StudentName: student.Name,
		Score:       student.Score,
	}

	return response, nil
}

func (repo *StudentRepoService) DeleteStudent(data schemas.DeleteStudent) error {
	student, err := repo.Repo.FindByID(data.StudentID)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return errors.New("student not found")
		}
		return errors.New("cannot get student")
	}

	err = mgm.Coll(student).Delete(student)
	if err != nil {
		return errors.New("cannot delete student")
	}

	return nil
}
