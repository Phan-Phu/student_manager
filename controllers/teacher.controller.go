package controllers

import (
	"net/http"
	"strings"
	"studenent_manager/models"
	"studenent_manager/models/db"
	"studenent_manager/models/repository"
	"studenent_manager/models/schemas"
	"studenent_manager/services"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateTeacher(c *gin.Context) {
	// Gán giá trị mặc định cho trường Role
	// cần check class id
	var requestBody schemas.RequestTeacher
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	teacher, err := services.CreateTeacher(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"teacher": teacher,
	}
	response.SendResponse(c)
}

func GetTeachers(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	teachers, err := services.GetTeachers()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"teachers": teachers,
	}
	response.SendResponse(c)
}

func GetTeacher(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	var requestBody struct {
		Name string `json:"name"`
	}

	_ = c.ShouldBindJSON(&requestBody)

	teacher, err := services.GetTeacher(requestBody.Name)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"teacher": teacher,
	}
	response.SendResponse(c)
}

func UpdateTeacher(c *gin.Context) {
	var requestBody *schemas.UpdateTeacher
	_ = c.ShouldBindJSON(&requestBody)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	teacher, err := services.UpdateTeacher(requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"teacher": teacher,
	}
	response.SendResponse(c)
}

func DeleteTeacher(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody struct {
		TeacherID int `json:"teacher_id"`
	}
	_ = c.ShouldBindJSON(&requestBody)

	err := services.DeleteTeacher(requestBody.TeacherID)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Message = "Teacher deleted successfully"
	response.SendResponse(c)
}

func LoginTeacher(c *gin.Context) {
	var requestBody *schemas.UpdateTeacher
	_ = c.ShouldBindJSON(&requestBody)

	teacher, err := repository.NewMongoTeacherRepository().FindByName(requestBody.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compare the hashed password
	isComparePassword := services.ComparePasswords(teacher.Password, requestBody.Password)
	if !isComparePassword {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	// Generate a JWT token
	session, err := services.GenerateJWTToken(teacher.ID.Hex(), requestBody.Username, db.TeacherRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save session to MongoDB
	err = mgm.Coll(session).Create(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	// Return token in response
	c.JSON(http.StatusOK, gin.H{
		"message":      "Logged in successfully",
		"tokenAccess":  session.AccessToken,
		"tokenRefresh": session.RefreshToken,
	})
}

func LogoutTeacher(c *gin.Context) {
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
