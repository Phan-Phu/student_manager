package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update with the correct import path for your db package

func TeacherMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID, exists := c.Get("teacherID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("teacherID", teacherID)
		c.Next()
	}
}
