package controllers

import (
	"net/http"
	"student_manager/models"
	"student_manager/models/schemas"
	"student_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateSubject(c *gin.Context) {
	var requestBody schemas.SubjectRequest
	_ = c.ShouldBindJSON(&requestBody)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	class, err := services.CreateSubject(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"subject": class,
	}
	response.SendResponse(c)
}
