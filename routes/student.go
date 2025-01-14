package routes

import (
	"studenent_manager/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoute(router *gin.RouterGroup) {
	auth := router.Group("/student")
	{
		auth.POST(
			"/create",
			controllers.CreateStudent,
		)

		auth.POST(
			"/getAll", //
			controllers.GetStudents,
		)

		auth.GET(
			"/getById",
			controllers.GetStudent,
		)

		auth.PUT(
			"/update",
			controllers.UpdateStudent,
		)

		auth.PATCH(
			"/updateScore",
			controllers.UpdateScore,
		)

		auth.DELETE(
			"/deleteById",
			controllers.DeleteStudent,
		)
	}
}
