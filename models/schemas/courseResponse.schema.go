package schemas

import (
	"student_manager/models/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SemesterResponse struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	CreateDate string           `json:"create_date"`
	EndDate    string           `json:"end_date"`
	Classes    []CourseResponse `json:"classes"`
}

type CourseResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Students  []db.Student         `json:"students"`
	Teachers  []db.Teacher         `json:"teachers"`
	Classes   ClassResponse        `json:"class"`
	Subject   SubjectResponse      `json:"subject"`
	StartDate string               `json:"start_date"`
	EndDate   string               `json:"end_date"`
	ScoreIDs  []primitive.ObjectID `json:"score_ids"`
}

type CourseDetail struct {
	Students []db.Student    `bson:"students" json:"students"`
	Teachers []db.Teacher    `bson:"teachers" json:"teachers"`
	Status   int             `bson:"status" json:"status"`
	Class    ClassResponse   `bson:"class" json:"class"`
	Subject  SubjectResponse `bson:"subject" json:"subject"`
}

type ScoreResponse struct {
	ScoreId string               `bson:"_id" json:"_id"`
	Student StudentDetails       `bson:"student" json:"student"`
	Subject SubjectResponse      `bson:"subject" json:"subject"`
	Score   map[db.ScoreType]int `bson:"score" json:"score"`
}

type ScoreDetails struct {
	ScoreId string               `bson:"_id" json:"_id"`
	Student StudentDetails       `bson:"student" json:"student"`
	Subject SubjectResponse      `bson:"subject" json:"subject"`
	Class   ClassResponse        `bson:"class" json:"class"`
	Score   map[db.ScoreType]int `bson:"score" json:"score"`
}
