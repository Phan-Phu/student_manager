package controllers

import (
	"net/http"
	"student_manager/models"
	"student_manager/models/schemas"
	"student_manager/services"

	"github.com/gin-gonic/gin"
)

func CreateCourse(ctx *gin.Context) {
	var req schemas.CreateCourseRequest
	_ = ctx.ShouldBindJSON(&req)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	err := req.Validate()
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(ctx)
		return
	}

	data, err := services.CourseService.CreateCourse(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"Data": data,
	}
	response.SendResponse(ctx)
}

func UpdateCourse(c *gin.Context) {
	var requestBody schemas.UpdateCourseRequest
	_ = c.ShouldBindJSON(&requestBody)
	id := c.Param("id")

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	score, err := services.CourseService.UpdateCourse(id, requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"course": score,
	}
	response.SendResponse(c)
}

func UpdateScore(c *gin.Context) {
	var requestBody schemas.UpdateScoreRequest
	_ = c.ShouldBindJSON(&requestBody)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	score, err := services.CourseService.UpdateScore(requestBody)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"score": score,
	}
	response.SendResponse(c)
}

func GetCourses(c *gin.Context) {
	var requestBody schemas.UpdateScoreRequest
	_ = c.ShouldBindJSON(&requestBody)

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
	}

	score, err := services.CourseService.GetCourses()
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Success = false
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.Data = gin.H{
		"score": score,
	}
	response.SendResponse(c)
}
