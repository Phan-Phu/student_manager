package controllers

import (
	"net/http"
	"strings"
	"studenent_manager/models"
	"studenent_manager/models/db"
	"studenent_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateClass(c *gin.Context) {
	var requestBody db.Class
	_ = c.ShouldBindJSON(&requestBody)
	requestBody.Name = strings.TrimSpace(requestBody.Name)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	class, err := services.CreateClass(requestBody.Name)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"class": class,
	}
	response.SendResponse(c)
}

func GetClasses(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	classes, err := services.GetClasses()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"classes": classes,
	}
	response.SendResponse(c)
}

func GetClass(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	var requestBody struct {
		Name string `json:"name"`
	}
	_ = c.ShouldBindJSON(&requestBody)

	class, err := services.GetClass(requestBody.Name)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"class": class,
	}
	response.SendResponse(c)
}

func UpdateClass(c *gin.Context) {
	// var requestBody db.Class
	// _ = c.ShouldBindJSON(&requestBody)
	// requestBody.Name = strings.TrimSpace(requestBody.Name)

	// response := &models.Response{
	// 	StatusCode: http.StatusOK,
	// 	Success:    true,
	// }

	// students := requestBody.StudentIds
	// teachers := requestBody.TeacherIds

	// class, err := services.UpdateClass(requestBody.ClassID, requestBody.Name, students, teachers)
	// if err != nil {
	// 	response.StatusCode = http.StatusNotFound
	// 	response.Success = false
	// 	response.Message = err.Error()
	// 	response.SendResponse(c)
	// 	return
	// }

	// response.Data = gin.H{
	// 	"class": class,
	// }
	// response.SendResponse(c)
}

func DeleteClass(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}
	var requestBody struct {
		ClassID int `json:"class_id"`
	}
	_ = c.ShouldBindJSON(&requestBody)

	err := services.DeleteClass(requestBody.ClassID)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Message = "Class deleted successfully"
	response.SendResponse(c)
}
