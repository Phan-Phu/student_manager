package schemas

import (
	"errors"
	"student_manager/models/db"
	"time"
)

// type CreateSemesterRequest struct {
// 	Name       string `json:"name"`
// 	CreateDate string `json:"create_date"`
// 	EndDate    string `json:"end_date"`
// 	ClassId    string `json:"class_id"` // not yet
// }

// func (r *CreateSemesterRequest) Validate() error {
// 	if r.Name == "" {
// 		return errors.New("name is required")
// 	}

// 	if _, err := time.Parse("2006-01-02", r.CreateDate); err != nil {
// 		return errors.New("invalid create_date format, expected YYYY-MM-DD")
// 	}
// 	if _, err := time.Parse("2006-01-02", r.EndDate); err != nil {
// 		return errors.New("invalid end_date format, expected YYYY-MM-DD")
// 	}

// 	return nil
// }

// type UpdateSemesterRequest struct {
// 	ID         string `json:"semester_id"`
// 	Name       string `json:"name"`
// 	CreateDate string `json:"create_date"`
// 	EndDate    string `json:"end_date"`
// }

// func (r *UpdateSemesterRequest) Validate() error {
// 	if r.Name == "" {
// 		return errors.New("name is required")
// 	}

// 	if _, err := time.Parse("2006-01-02", r.CreateDate); err != nil {
// 		return errors.New("invalid create_date format, expected YYYY-MM-DD")
// 	}
// 	if _, err := time.Parse("2006-01-02", r.EndDate); err != nil {
// 		return errors.New("invalid end_date format, expected YYYY-MM-DD")
// 	}

// 	return nil
// }

type CreateCourseRequest struct {
	Name         string   `json:"name"`
	StudentIds   []string `json:"student_ids"`
	TeacherIds   []string `json:"teacher_ids"`
	ClassId      string   `json:"class_id"`
	SubjectId    string   `json:"subject_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Examinations []int    `json:"examinations"`
}

type UpdateCourseRequest struct {
	StudentIds []string `json:"student_ids"`
	TeacherIds []string `json:"teacher_ids"`
	ClassId    string   `json:"class_id"`
	SubjectId  string   `json:"subject_id"`
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
}

const dateFormat = "2006-01-02"

func (req *CreateCourseRequest) Validate() error {
	if len(req.StudentIds) == 0 {
		return errors.New("student_ids cannot be empty")
	}

	if len(req.TeacherIds) == 0 {
		return errors.New("teacher_ids cannot be empty")
	}

	if req.ClassId == "" {
		return errors.New("class_id cannot be empty")
	}

	if req.SubjectId == "" {
		return errors.New("subject_id cannot be empty")
	}

	if req.StartDate == "" {
		return errors.New("start_date cannot be empty")
	}
	if _, err := time.Parse(dateFormat, req.StartDate); err != nil {
		return errors.New("start_date must be in the format YYYY-MM-DD")
	}

	if req.EndDate == "" {
		return errors.New("end_date cannot be empty")
	}
	if _, err := time.Parse(dateFormat, req.EndDate); err != nil {
		return errors.New("end_date must be in the format YYYY-MM-DD")
	}

	startDate, _ := time.Parse(dateFormat, req.StartDate)
	endDate, _ := time.Parse(dateFormat, req.EndDate)
	if endDate.Before(startDate) {
		return errors.New("end_date must be after start_date")
	}

	// Check if Examinations is not empty
	if len(req.Examinations) == 0 {
		return errors.New("examinations cannot be empty")
	}

	// Ensure that each examination type in Examinations is valid
	for _, exam := range req.Examinations {
		if exam != int(db.SemiSemester) &&
			exam != int(db.FinalSemester) &&
			exam != int(db.Presentation) {
			return errors.New("invalid examination type in examinations")
		}
	}

	return nil
}

type UpdateScoreRequest struct {
	ScoreId string `json:"score_id"`
	Type    int    `json:"type"`
	Score   int    `json:"score"`
}
