package services

import (
	"errors"
	db "studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStudent(studentId int, name string, classID int, birthDay string) (*db.Student, error) {
	student := db.CreateStudent(studentId, name, classID, birthDay, 0)
	err := mgm.Coll(student).Create(student)

	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return student, nil
}

func GetStudents() ([]db.Student, error) {
	var students []db.Student
	err := mgm.Coll(&db.Student{}).SimpleFind(&students, bson.M{})
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

func UpdateStudent(studentId int, name string, classID int, birthDay string) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"student_id": studentId}, student)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	student.Name = name
	student.ClassID = classID
	student.BirthDay = birthDay

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New("cannot update student")
	}

	return student, nil
}

func UpdateScore(studentId int, score int) (*db.Student, error) {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"student_id": studentId}, student)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("cannot get student")
	}

	student.Score = score

	err = mgm.Coll(student).Update(student)
	if err != nil {
		return nil, errors.New("cannot update student")
	}

	return student, nil
}

func DeleteStudent(studentID int) error {
	student := &db.Student{}
	err := mgm.Coll(student).First(bson.M{"student_id": studentID}, student)
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
