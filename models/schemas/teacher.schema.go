package schemas

import (
	"errors"
	"strings"
	"time"
)

type RequestTeacher struct {
	Name     string `json:"name"`
	BirthDay string `json:"birth_day"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate RequestTeacher
func (r *RequestTeacher) Validate() error {
	// Check name
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	// if len(r.Name) > 50 {
	// 	return errors.New("name is too long")
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

type UpdateTeacher struct {
	TeacherID int    `json:"teacher_id"`
	Name      string `json:"name"`
	ClassIDs  []int  `json:"class_id"`
	BirthDay  string `json:"birth_day"`
	Score     int    `json:"score"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`
}

func (u *UpdateTeacher) Validate() error {
	if u.TeacherID <= 0 {
		return errors.New("invalid teacher_id")
	}

	_, err := time.Parse("02-01-2006", u.BirthDay)
	if err != nil {
		return errors.New("invalid birth_day format, should be DD-MM-YYYY")
	}

	if u.Score < 0 || u.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	return nil
}

type DeleteTeacher struct {
	TeacherID int `json:"teacher_id"`
}

func (d *DeleteTeacher) Validate() error {
	if d.TeacherID <= 0 {
		return errors.New("invalid teacher_id")
	}
	return nil
}
