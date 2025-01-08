package routes

import (
	"studenent_manager/controllers"
	"studenent_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func ClassRoute(router *gin.RouterGroup) {
	auth := router.Group("/class")
	{
		auth.POST(
			"/create",
			middlewares.CreateClass(),
			controllers.CreateClass,
		)

		auth.GET(
			"/getAll",
			controllers.GetClass,
		)

		auth.GET(
			"/getById",
			controllers.GetClass,
		)

		auth.POST(
			"/update",
			middlewares.CreateClass(),
			controllers.UpdateClass,
		)

		auth.DELETE(
			"/deleteById",
			middlewares.DeleteClass(),
			controllers.DeleteClass,
		)
	}
}
