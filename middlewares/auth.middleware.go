package middlewares

import (
	"net/http"
	"studenent_manager/models/db"
	"studenent_manager/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// Middleware để check JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		session := &db.Token{}
		err := mgm.Coll(session).First(bson.M{
			"access_token": accessToken,
			"access_expires_at": bson.M{
				"$gt": time.Now(),
			},
		}, session)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error token not exist"})
			c.Abort()
			return
		}

		response, errToken := services.ValidateJWTToken(accessToken)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error validating JWT token"})
			c.Abort()
			return
		}

		// Set claims to context
		c.Set("username", response.Username)
		c.Set("userId", response.UserId)
		c.Next()
	}
}

// func CleanupExpiredSessions() {
// 	_, err := mgm.Coll(&db.Token{}).DeleteMany(mgm.Ctx(), bson.M{
// 		"refresh_expires_at": bson.M{
// 			"$lt": time.Now(),
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("Failed to cleanup expired sessions: %v", err)
// 	}
// }
