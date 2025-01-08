package services

import (
	"errors"
	db "studenent_manager/models/db"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateClass(name string, studentIds []int, teacherIds []int) (*db.Class, error) {
	classId := GenerateClassID()
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	class := db.CreateClass(classId, name, currentTime, studentIds, teacherIds)
	err := mgm.Coll(class).Create(class)

	if err != nil {
		return nil, errors.New("cannot create new class")
	}

	class, _ = loadNavigationProperty(class)
	return class, nil
}

func GetClasses() ([]db.Class, error) {
	var classes []db.Class
	err := mgm.Coll(&db.Class{}).SimpleFind(&classes, bson.M{})
	if err != nil {
		return nil, errors.New("cannot get classes")
	}
	for i := 0; i < len(classes); i++ {

		class := &classes[i]
		class, _ = loadNavigationProperty(class)

		classes[i].Students = class.Students
		classes[i].Teachers = class.Teachers
	}

	return classes, nil
}

func GetClass(className string) (*db.Class, error) {
	class := &db.Class{}
	err := mgm.Coll(class).First(bson.M{"name": className}, class)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("class not found")
		}
		return nil, errors.New("cannot get class")
	}
	class, _ = loadNavigationProperty(class)

	return class, nil
}

func UpdateClass(classId int, name string, students []int, teachers []int) (*db.Class, error) {
	class := &db.Class{}
	err := mgm.Coll(class).First(bson.M{"class_id": classId}, class)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return nil, errors.New("class not found")
		}
		return nil, errors.New("cannot get class")
	}

	class.Name = name
	class.StudentIds = students
	class.TeacherIds = teachers

	err = mgm.Coll(class).Update(class)
	if err != nil {
		return nil, errors.New("cannot update class")
	}

	return class, nil
}

func DeleteClass(classID int) error {
	class := &db.Class{}
	err := mgm.Coll(class).First(bson.M{"class_id": classID}, class)
	if err != nil {
		if err == mgm.Ctx().Err() {
			return errors.New("class not found")
		}
		return errors.New("cannot get class")
	}

	err = mgm.Coll(class).Delete(class)
	if err != nil {
		return errors.New("cannot delete class")
	}

	return nil
}

func GenerateClassID() int {
	classes, _ := GetClasses()
	classId := len(classes) + 1
	return classId
}

func loadNavigationProperty(class *db.Class) (*db.Class, error) {
	var students []db.Student
	var teachers []db.Teacher

	// Fetch students if there are student IDs
	if len(class.StudentIds) > 0 {
		query := bson.M{"student_id": bson.M{"$in": class.StudentIds}}
		err := mgm.Coll(&db.Student{}).SimpleFind(&students, query)
		if err != nil {
			return class, errors.New("cannot find students")
		}
	}

	// Fetch teachers if there are teacher IDs
	if len(class.TeacherIds) > 0 {
		query := bson.M{"teacher_id": bson.M{"$in": class.TeacherIds}}
		err := mgm.Coll(&db.Teacher{}).SimpleFind(&teachers, query)
		if err != nil {
			return class, errors.New("cannot find teachers")
		}
	}

	class.Students = students
	class.Teachers = teachers

	return class, nil
}
