package services

import (
	db "studenent_manager/models/db"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost là chi phí mặc định (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func GenerateJWTToken(userName string, role db.Role) (*db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	accessClaims := &db.Claims{
		Username: userName,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpiresAt.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour)
	refreshClaims := &db.Claims{
		Username: userName,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	session := db.NewToken(accessTokenString, refreshTokenString, db.AdminRole, accessExpiresAt, refreshExpiresAt)

	return session, nil
}
