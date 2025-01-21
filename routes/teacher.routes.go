package routes

import (
	"student_manager/controllers"

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

		auth.PATCH(
			"/updateInfo",
			controllers.UpdateInfoTeacher,
		)

		auth.PATCH(
			"/updateClass",
			controllers.UpdateClassTeacher,
		)

		auth.DELETE(
			"/deleteById",
			controllers.DeleteTeacher,
		)
	}
}
