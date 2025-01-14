package db

import (
	"errors"
	"regexp"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID       primitive.ObjectID `json:"_id"`      // unique
	Username string             `json:"username"` // unique
	Password string             `json:"password"`
	Role     Role               `json:"role"`
}

type RequestAccount struct {
	Username string `json:"username"` // unique
	Password string `json:"password"`
}

func (a *RequestAccount) Validate() error {
	// Validate username (must be non-empty)
	if len(a.Username) == 0 {
		return errors.New("username cannot be empty")
	}

	// Validate password (must be at least 8 characters and contain at least one number and one letter)
	if len(a.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	match, _ := regexp.MatchString(`[A-Za-z]`, a.Password)
	if !match {
		return errors.New("password must contain at least one letter")
	}

	match, _ = regexp.MatchString(`[0-9]`, a.Password)
	if !match {
		return errors.New("password must contain at least one number")
	}

	return nil
}

type ResponseAccount struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.StandardClaims
}
