package controllers

import (
	"net/http"
	"studenent_manager/models/db"
	"studenent_manager/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var adminAccount = db.Account{
	ID:       primitive.NewObjectID(),
	Username: "admin",
	Password: "admin123",
	Role:     db.AdminRole,
}

func LoginAdmin(c *gin.Context) {
	var account db.RequestAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials
	if account.Username != adminAccount.Username || account.Password != adminAccount.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	data, err := services.GenerateJWTToken(adminAccount.ID.Hex(), account.Username, db.AdminRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return token in responses
	c.JSON(http.StatusOK, gin.H{
		"message":      "Logged in successfully",
		"tokenAccess":  data.AccessToken,
		"tokenRefresh": data.RefreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	// Get refresh token from request
	refreshTokenString := c.GetHeader("refresh-Token")
	if refreshTokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Validate refresh token
	claims := &db.ResponseAccount{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("refresh-token"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Find session in database
	session := &db.Token{}
	err = mgm.Coll(session).First(bson.M{
		"refresh_token": refreshTokenString,
		"refresh_expires_at": bson.M{
			"$gt": time.Now(),
		},
	}, session)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Create new access token
	accessExpiresAt := time.Now().Add(30 * time.Minute)
	accessClaims := &db.ResponseAccount{
		UserId:   claims.UserId,
		Username: claims.Username,
		Role:     claims.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpiresAt.Unix(),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	newAccessTokenString, err := newAccessToken.SignedString([]byte("admin-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update session in database
	session.AccessToken = newAccessTokenString
	session.AccessExpiresAt = accessExpiresAt

	err = mgm.Coll(session).Update(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessTokenString,
	})
}

func LogoutAdmin(c *gin.Context) {
	// Get token from Authorization header
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	// Delete session from MongoDB
	result, err := mgm.Coll(&db.Token{}).DeleteOne(mgm.Ctx(), bson.M{
		"access_token": accessToken,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
