package controllers

import (
	"net/http"
	"studenent_manager/models/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var adminAccount = db.Account{
	Username: "admin",
	Password: "admin123",
	Role:     db.AdminRole,
}

func Login(c *gin.Context) {
	var account db.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials
	if account.Username != adminAccount.Username || account.Password != adminAccount.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &db.Claims{
		Username: account.Username,
		Role:     account.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the token in a cookie
	c.SetCookie("token", tokenString, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Logout(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
