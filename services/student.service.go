package services

import (
	"errors"
	"fmt"
	"math"
	db "studenent_manager/models/db"
	"studenent_manager/models/schemas"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateStudent(data schemas.RequestStudent) (*db.Student, error) {
	studentId := GenerateStudentID()
	initScore := 0

	student := db.NewStudent(studentId, data.Name, data.ClassID, data.BirthDay, initScore)
	err := StudentRepository.Create(student)

	if err != nil {
		return nil, errors.New("cannot create new student")
	}

	return student, nil
}

func GetStudents() ([]*db.Student, error) {
	students, err := StudentRepository.FindAll()

	if err != nil {
		return nil, errors.New("cannot get students")
	}
	return students, nil
}

func GetStudentsWithPagination(data *schemas.PaginationRequest) (*schemas.PaginationResponse[db.Student], error) {
	findOptions := options.Find()
	findOptions.SetSkip((int64(data.Page) - 1) * data.Limit)
	findOptions.SetLimit(int64(data.Limit))

	students, _ := StudentRepository.FindByFindOptions(findOptions)
	total, _ := StudentRepository.Count()

	response := &schemas.PaginationResponse[db.Student]{
		Page:      data.Page,
		Limit:     data.Limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(data.Limit))),
		Data:      students,
	}
	return response, nil
}

func GetStudent(studentName string) (*db.Student, error) {
	start := time.Now()
	student, err := StudentRepository.FindByName(studentName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	elapsed := time.Since(start)
	fmt.Printf("Found student: %+v\n", student)
	fmt.Printf("Time with index: %v\n", elapsed)
	return student, nil
}

func UpdateStudent(data schemas.UpdateStudent) (*db.Student, error) {
	student, err := StudentRepository.FindByID(data.StudentID)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	student.Name = data.Name
	student.ClassID = data.ClassID
	student.BirthDay = data.BirthDay

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New("cannot update student")
	}

	return student, nil
}

func UpdateScore(data schemas.UpdateScoreStudent) (*db.Student, error) {
	student, err := StudentRepository.FindByID(data.StudentID)
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

	return student, nil
}

func DeleteStudent(data schemas.DeleteStudent) error {
	student, err := StudentRepository.FindByID(data.StudentID)
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

func GenerateStudentID() int {
	students, _ := GetStudents()
	studentId := len(students) + 1
	return studentId
}
