package schemas

import (
	"errors"
	"strings"
	"time"
)

type RequestStudent struct {
	Name     string `json:"name"`
	ClassId  string `json:"class_id"`
	Age      int    `json:"age"`
	BirthDay string `json:"birth_day"`
}

// Validate RequestStudent
func (r *RequestStudent) Validate() error {
	// Check name
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	// if len(r.Name) > 50 {
	// 	return errors.New("name is too long")
	// }

	// if r.ClassID <= 0 {
	// 	return errors.New("invalid class_id")
	// }
	if r.BirthDay == "" {
		return errors.New("birth_day is required")
	}

	_, err := time.Parse("02-01-2006", r.BirthDay)
	if err != nil {
		return errors.New("invalid birth_day format, should be DD-MM-YYYY")
	}

	return nil
}

type UpdateScoreStudent struct {
	StudentID string `json:"_id"`
	Score     int    `json:"score"`
}

type UpdateStudentRequest struct {
	Name     string `json:"name"`
	ClassID  string `json:"class_id"`
	BirthDay string `json:"birth_day"`
	Age      int    `json:"age"`
	Score    int    `json:"score"`
	IsLock   bool   `json:"isLock"`
}

func (u *UpdateStudentRequest) Validate() error {
	// if u.StudentID <= 0 {
	// 	return errors.New("invalid student_id")
	// }

	_, err := time.Parse("02-01-2006", u.BirthDay)
	if err != nil {
		return errors.New("invalid birth_day format, should be DD-MM-YYYY")
	}

	if u.Score < 0 || u.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	return nil
}

func (u *UpdateScoreStudent) Validate() error {
	// if u.StudentID <= 0 {
	// 	return errors.New("invalid student_id")
	// }

	if u.Score < 0 || u.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	return nil
}

type DeleteStudent struct {
	StudentID string `json:"_id"`
}

func (d *DeleteStudent) Validate() error {
	// if d.StudentID <= 0 {
	// 	return errors.New("invalid student_id")
	// }
	return nil
}
