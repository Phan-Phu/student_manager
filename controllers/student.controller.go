package controllers

import (
	"net/http"
	"strings"
	"studenent_manager/models"
	"studenent_manager/models/schemas"
	"studenent_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody schemas.RequestStudent
	_ = c.ShouldBindJSON(&requestBody)

	err := requestBody.Validate()
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	student, err := services.CreateStudent(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"student": student,
	}
	response.SendResponse(c)
}

func GetStudents(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	students, err := services.GetStudents()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"students": students,
	}
	response.SendResponse(c)
}

func GetStudent(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	var requestBody struct {
		Name string `json:"name"`
	}
	_ = c.ShouldBindJSON(&requestBody)

	student, err := services.GetStudent(requestBody.Name)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"student": student,
	}
	response.SendResponse(c)
}

func UpdateStudent(c *gin.Context) {
	var requestBody schemas.UpdateStudent
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	student, err := services.UpdateStudent(requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"student": student,
	}
	response.SendResponse(c)
}

func UpdateScore(c *gin.Context) {
	var requestBody schemas.UpdateScoreStudent
	_ = c.ShouldBindJSON(&requestBody)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	student, err := services.UpdateScore(requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"student": student,
	}
	response.SendResponse(c)
}

func DeleteStudent(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}
	var requestBody struct {
		StudentID int `json:"student_id"`
	}
	_ = c.ShouldBindJSON(&requestBody)

	err := services.DeleteStudent(requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Message = "Student deleted successfully"
	response.SendResponse(c)
}
