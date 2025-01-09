package services

import (
	"errors"
	db "studenent_manager/models/db"
	"studenent_manager/models/repository"
	"studenent_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStudent(data schemas.RequestStudent) (*db.Student, error) {
	studentId := GenerateStudentID()
	initScore := 0

	student := db.CreateStudent(studentId, data.Name, data.ClassID, data.BirthDay, initScore)
	err := repository.NewMongoStudentRepository().Create(student)

	if err != nil {
		return nil, errors.New("cannot create new student")
	}

	return student, nil
}

func GetStudents() ([]*db.Student, error) {
	students, err := repository.NewMongoStudentRepository().FindAll()

	if err != nil {
		return nil, errors.New("cannot get students")
	}
	return students, nil
}

func GetStudent(studentName string) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"name": studentName}, student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}
	return student, nil
}

func UpdateStudent(data schemas.UpdateStudent) (*db.Student, error) {
	student, err := repository.NewMongoStudentRepository().FindByID(data.StudentID)
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
	student, err := repository.NewMongoStudentRepository().FindByID(data.StudentID)
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
	student, err := repository.NewMongoStudentRepository().FindByID(data.StudentID)
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
