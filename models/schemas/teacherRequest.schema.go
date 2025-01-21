package schemas

import (
	"errors"
	"strings"
	"time"
)

type RequestTeacher struct {
	Name     string   `json:"name"`
	ClassIds []string `json:"class_ids"`
	BirthDay string   `json:"birth_day"`
	Username string   `json:"userName"`
	Password string   `json:"password"`
}

func (r *RequestTeacher) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}

	if r.BirthDay == "" {
		return errors.New("birth_day is required")
	}

	_, err := time.Parse("02-01-2006", r.BirthDay)
	if err != nil {
		return errors.New("invalid birth_day format, should be DD-MM-YYYY")
	}

	return nil
}

type UpdateInfoTeacherRequest struct {
	TeacherID string `json:"teacher_id"`
	Name      string `json:"name"`
	BirthDay  string `json:"birth_day"`
}

func (u *UpdateInfoTeacherRequest) Validate() error {
	_, err := time.Parse("02-01-2006", u.BirthDay)
	if err != nil {
		return errors.New("invalid birth_day format, should be DD-MM-YYYY")
	}

	return nil
}

type UpdateTeacherClassRequest struct {
	TeacherID     string          `json:"teacher_id"`
	ClassesUpdate []ClassIDUpdate `json:"classes_update"`
}

type ClassIDUpdate struct {
	OldClassID string `json:"old_class_id"`
	NewClassID string `json:"new_class_id"`
}

type DeleteTeacher struct {
	TeacherID string `json:"teacher_id"`
}

type LoginTeacherRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
