package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"studenent_manager/models"
	"studenent_manager/models/db"
	"studenent_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	// Gán giá trị mặc định cho trường Role
	var requestBody db.Student
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// class := GetClassById(requestBody.ClassID)
	// studentId := len(class.Students) + 1
	studentId := 1

	student, err := services.CreateStudent(studentId, requestBody.Name, requestBody.ClassID, requestBody.BirthDay)
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
		StatusCode: http.StatusOK,
		Success:    true,
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
	var requestBody db.Student
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	student, err := services.UpdateStudent(requestBody.StudentID, requestBody.Name, requestBody.ClassID, requestBody.BirthDay)
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

	studentIDStr := c.Param("student_id")
	studentID, err := strconv.Atoi(studentIDStr)

	err = services.DeleteStudent(studentID)
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
