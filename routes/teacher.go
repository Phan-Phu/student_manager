package routes

import (
	"studenent_manager/controllers"

	"github.com/gin-gonic/gin"
)

func TeacherRoute(router *gin.RouterGroup) {
	auth := router.Group("/teacher")
	{
		auth.POST(
			"/create",
			controllers.CreateTeacher,
		)

		auth.GET(
			"/getAll",
			controllers.GetTeachers,
		)

		auth.GET(
			"/getById",
			controllers.GetTeacher,
		)

		auth.POST(
			"/update",
			controllers.UpdateTeacher,
		)

		auth.DELETE(
			"/deleteById",
			controllers.DeleteTeacher,
		)
	}
}
