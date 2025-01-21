package routes

import (
	"student_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SubjectRoute(router *gin.RouterGroup) {
	auth := router.Group("/subject")
	{
		auth.POST(
			"/create",
			controllers.CreateSubject,
		)

		auth.PATCH(
			"/update",
			// controllers.UpdateSemester,
		)

	}
}
