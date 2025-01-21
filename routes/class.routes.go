package routes

import (
	"student_manager/controllers"

	"github.com/gin-gonic/gin"
)

func ClassRoute(router *gin.RouterGroup) {

	router.GET(
		"/classes",
		controllers.GetClasses,
	)

	router.POST(
		"/classes",
		controllers.CreateClass,
	)

	router.GET(
		"/classes/:id",
		controllers.GetClass,
	)

	router.PUT(
		"/classes/:id",
		controllers.UpdateClass,
	)

	router.DELETE(
		"/classes/:id",
		controllers.DeleteClass,
	)
}
