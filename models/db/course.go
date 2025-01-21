package db

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreType int

const (
	FinalSemester ScoreType = 0
	SemiSemester  ScoreType = 1
	Presentation  ScoreType = 2
	// etc...
)

type Course struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string               `bson:"name"`
	ClassId          primitive.ObjectID   `bson:"class_id"`
	TeacherIds       []primitive.ObjectID `bson:"teacher_ids"`
	StudentIds       []primitive.ObjectID `bson:"student_ids"`
	SubjectId        primitive.ObjectID   `bson:"subject_id"`
	StartDate        string               `bson:"create_date"`
	EndDate          string               `bson:"end_date"`
	Scores           []primitive.ObjectID `bson:"score_ids"`
}

type Score struct {
	mgm.DefaultModel `bson:",inline"`
	StudentId        primitive.ObjectID `bson:"student_id"`
	SubjectId        primitive.ObjectID `bson:"subject_id"`
	Score            map[ScoreType]int  `bson:"score"`
}

func NewCourse(
	name string,
	studentIds []primitive.ObjectID,
	teacherIds []primitive.ObjectID,
	classId primitive.ObjectID,
	subjectId primitive.ObjectID,
	startDate string,
	endDate string,
	scoreIds []primitive.ObjectID) *Course {
	return &Course{
		Name:       name,
		StudentIds: studentIds,
		TeacherIds: teacherIds,
		ClassId:    classId,
		SubjectId:  subjectId,
		StartDate:  startDate,
		EndDate:    endDate,
		Scores:     scoreIds,
	}
}

func NewScore(
	studentId primitive.ObjectID,
	subjectId primitive.ObjectID,
	score map[ScoreType]int) *Score {
	return &Score{
		StudentId: studentId,
		SubjectId: subjectId,
		Score:     score,
	}
}
