package services

import (
	"context"
	"errors"
	"log"
	"student_manager/models"
	db "student_manager/models/db"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GenerateJWTToken(userId string, userName string, role db.Role) (*db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	accessClaims := &db.ResponseAccount{
		UserId:   userId,
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
	refreshClaims := &db.ResponseAccount{
		UserId:   accessClaims.UserId,
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

	refreshIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "refresh_expires_at", Value: "hashed"}},
		Options: options.Index().SetExpireAfterSeconds(int32(Config.JWTRefreshExpirationDays)),
	}

	_, err = mgm.Coll(&db.Token{}).Indexes().CreateMany(context.Background(), []mongo.IndexModel{refreshIndexModel})
	if err != nil {
		log.Fatal(err)
	}

	// Save session to MongoDB
	err = mgm.Coll(session).Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func ValidateJWTToken(token string) (*db.ResponseAccount, error) {
	response := &db.ResponseAccount{}
	_, err := jwt.ParseWithClaims(token, response, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})
	return response, err
}

func ConvertStringsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	results := []primitive.ObjectID{}
	for i := 0; i < len(ids); i++ {
		id, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
		}
		results = append(results, id)
	}
	return results, nil
}

func ConvertStringToObjectID(id string) (*primitive.ObjectID, error) {
	result, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(models.GetErrorMessage(models.ErrorCodeIDIsWrong))
	}
	return &result, nil
}
