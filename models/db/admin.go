package db

import "github.com/golang-jwt/jwt"

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.StandardClaims
}
