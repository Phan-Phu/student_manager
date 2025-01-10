package services

import (
	"errors"
	"log"
	"strings"
	db "studenent_manager/models/db"
	"studenent_manager/models/schemas"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTeacher(data schemas.RequestTeacher) (*db.Teacher, error) {

	teacherId := GenerateTeacherID()

	isUserNameExist := checkUserNameTeacher(data.Username)
	if !isUserNameExist {
		return nil, errors.New("UserName is exist")
	}

	hashPassword, err := HashPassword(data.Password)
	if err != nil {
		return nil, errors.New("cannot hash password")
	}

	teacher := db.CreateTeacher(teacherId, data.Name, nil, data.BirthDay, data.Username, hashPassword)

	err = mgm.Coll(teacher).Create(teacher)

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

func UpdateTeacher(data *schemas.UpdateTeacher) (*db.Teacher, error) {
	teacher := &db.Teacher{}
	teacherName := strings.TrimSpace(data.Name)

	err := mgm.Coll(teacher).First(bson.M{"teacher_id": data.TeacherID}, teacher)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("cannot get teacher")
	}

	teacher.Name = teacherName
	teacher.ClassIds = data.ClassIDs
	teacher.BirthDay = data.BirthDay

	err = mgm.Coll(teacher).Update(teacher)
	if err != nil {
		return nil, errors.New("cannot update teacher")
	}

	updatedTeacher := &db.Teacher{}
	err = mgm.Coll(updatedTeacher).First(bson.M{"teacher_id": data.TeacherID}, updatedTeacher)
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

func checkUserNameTeacher(username string) bool {
	// Get the default collection for the User model
	collection := mgm.Coll(&db.Teacher{})

	// Count documents with the given username
	filter := bson.M{"user_name": username}
	count, err := collection.CountDocuments(mgm.Ctx(), filter, options.Count())
	if err != nil {
		return false
	}

	return count == 0
}
