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

		auth.GET(
			"/getAll",
			controllers.GetStudents,
		)

		auth.GET(
			"/getById",
			controllers.GetStudent,
		)

		auth.POST(
			"/update",
			controllers.UpdateStudent,
		)

		auth.POST(
			"/updateScore",
			controllers.UpdateScore,
		)

		auth.DELETE(
			"/deleteById",
			controllers.DeleteStudent,
		)
	}
}
