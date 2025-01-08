package db

type Role int

const (
	AdminRole   Role = 0 // admin can be fix in first release.
	TeacherRole Role = 1
	StudentRole Role = 2
)

type Teacher struct {
	TeacherID   int     `json:"id" bson:"id"`
	TeacherName string  `json:"name" bson:"name"`
	BirthDay    float64 `json:"birth_day" bson:"birth_day"`
	Class       []int   `json:"class" bson:"class"`
	Role        Role    `json:"role" bson:"role"` // not input from user
}
