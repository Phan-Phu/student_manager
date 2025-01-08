package services

import (
	"errors"
	"log"
	db "studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTeacher(teacherId int, name string, classIds []int, birthDay string, userName string, password string) (*db.Teacher, error) {
	teacher := db.CreateTeacher(teacherId, name, classIds, birthDay, userName, password)
	err := mgm.Coll(teacher).Create(teacher)

	if err != nil {
		return nil, errors.New("cannot create new teacher")
	}

	return teacher, nil
}

func GetTeachers() ([]db.Teacher, error) {
	var teachers []db.Teacher
	err := mgm.Coll(&db.Teacher{}).SimpleFind(&teachers, bson.M{})
	if err != nil {
		return nil, errors.New("cannot get teachers")
	}
	return teachers, nil
}

func GetTeacher(teacherName string) (*db.Teacher, error) {
	teacher := &db.Teacher{}
	err := mgm.Coll(teacher).First(bson.M{"name": teacherName}, teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("cannot get teacher")
	}
	return teacher, nil
}

func UpdateTeacher(teacherId int, name string, classIds []int, birthDay string) (*db.Teacher, error) {
	teacher := &db.Teacher{}
	err := mgm.Coll(teacher).First(bson.M{"teacher_id": teacherId}, teacher)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("cannot get teacher")
	}

	teacher.Name = name
	teacher.ClassIds = classIds
	teacher.BirthDay = birthDay

	err = mgm.Coll(teacher).Update(teacher)
	if err != nil {
		return nil, errors.New("cannot update teacher")
	}

	updatedTeacher := &db.Teacher{}
	err = mgm.Coll(updatedTeacher).First(bson.M{"teacher_id": teacherId}, updatedTeacher)
	if err != nil {
		return nil, errors.New("cannot fetch updated teacher")
	}
	log.Println("Updated Teacher:", updatedTeacher)

	return teacher, nil
}

func DeleteTeacher(teacherID int) error {
	teacher := &db.Teacher{}
	err := mgm.Coll(teacher).First(bson.M{"teacher_id": teacherID}, teacher)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return errors.New("teacher not found")
		}
		return errors.New("cannot get teacher")
	}

	err = mgm.Coll(teacher).Delete(teacher)
	if err != nil {
		return errors.New("cannot delete teacher")
	}

	return nil
}

func GenerateTeacherID() int {
	teachers, _ := GetTeachers()
	teacherId := len(teachers) + 1
	return teacherId
}

func LoginTeacher(userName string, password string) error {
	var teacher db.Teacher
	err := mgm.Coll(&teacher).First(bson.M{"username": userName}, &teacher)
	if err != nil {
		return errors.New("cannot login teacher")
	}

	if teacher.Password != password {
		return errors.New("cannot login teacher")
	}

	return nil
}
