package routes

import (
	"student_manager/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoute(router *gin.RouterGroup) {

	router.POST(
		"/students",
		controllers.CreateStudent,
	)

	router.GET(
		"/students",
		controllers.GetStudents,
	)

	router.GET(
		"/students/:id",
		controllers.GetStudent,
	)

	router.PUT(
		"/students/:id",
		controllers.UpdateStudent,
	)

	router.PATCH(
		"/students/:id",
		controllers.UpdateScore,
	)

	router.DELETE(
		"/controllers/:id",
		controllers.DeleteStudent,
	)

}
