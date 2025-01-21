package controllers

// import (
// 	"net/http"
// 	"student_manager/models"
// 	"student_manager/models/schemas"
// 	"student_manager/services"

// 	"github.com/gin-gonic/gin"
// )

// // type SemesterController struct {
// // 	service *services.SemesterRepoService
// // }

// // func NewSemesterController(service *services.SemesterRepoService) *SemesterController {
// // 	return &SemesterController{service: service}
// // }

// func CreateSemester(ctx *gin.Context) {
// 	var req schemas.CreateSemesterRequest
// 	_ = ctx.ShouldBindJSON(&req)

// 	response := &models.Response{
// 		StatusCode: http.StatusBadRequest,
// 		Success:    false,
// 	}

// 	err := req.Validate()
// 	if err != nil {
// 		response.Message = err.Error()
// 		response.SendResponse(ctx)
// 		return
// 	}

// 	data, err := services.SemesterService.CreateSemester(req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response.StatusCode = http.StatusOK
// 	response.Success = true
// 	response.Data = gin.H{
// 		"Data": data,
// 	}
// 	response.SendResponse(ctx)
// }

// func UpdateSemester(ctx *gin.Context) {
// 	// var req schemas.UpdateSemester
// 	// _ = c.ShouldBindJSON(&requestBody)

// 	// err := req.Validate()
// 	// if err != nil {
// 	// 	response.Message = err.Error()
// 	// 	response.SendResponse(c)
// 	// 	return
// 	// }
// 	// response := &models.Response{
// 	// 	StatusCode: http.StatusBadRequest,
// 	// 	Success:    false,
// 	// }

// 	// response, err := c.service.UpdateSemester(req)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// response.StatusCode = http.StatusOK
// 	// response.Success = true
// 	// response.Data = response
// 	// response.SendResponse(c)
// }
