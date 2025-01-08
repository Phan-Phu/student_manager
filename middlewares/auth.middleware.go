package middlewares

import (
	"net/http"
	"studenent_manager/models/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Update with the correct import path for your db package

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the cookies
		cookie, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort() // Prevent further processing
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort() // Prevent further processing
			return
		}

		tokenString := cookie
		claims := &db.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort() // Prevent further processing
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort() // Prevent further processing
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Prevent further processing
			return
		}

		// Check if the user has the correct role
		if claims.Role != db.AdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort() // Prevent further processing
			return
		}

		// Continue with the next handler in the chain
		c.Next()
	}
}
