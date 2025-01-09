package controllers

import (
	"net/http"
	"strings"
	"studenent_manager/models"
	"studenent_manager/models/db"
	"studenent_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateTeacher(c *gin.Context) {
	// Gán giá trị mặc định cho trường Role
	// cần check class id
	var requestBody db.Teacher
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	teacherId := services.GenerateTeacherID()

	teacher, err := services.CreateTeacher(teacherId, requestBody.Name, requestBody.ClassIds,
		requestBody.BirthDay, requestBody.Username, requestBody.Password)
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
		StatusCode: http.StatusOK,
		Success:    true,
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
	var requestBody db.Teacher
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	teacher, err := services.UpdateTeacher(requestBody.TeacherID, requestBody.Name, requestBody.ClassIds, requestBody.BirthDay)
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
		StatusCode: http.StatusOK,
		Success:    true,
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
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errLogin := services.LoginTeacher(requestBody.Username, requestBody.Password)
	if errLogin != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Login Fail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func LogoutTeacher(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
