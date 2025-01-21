package routes

import (
	"student_manager/controllers"

	"github.com/gin-gonic/gin"
)

func CourseRoute(router *gin.RouterGroup) {
	router.POST(
		"/courses",
		controllers.CreateCourse,
	)

	router.PATCH(
		"/courses/",
		controllers.UpdateScore,
	)

	router.PATCH(
		"/courses/score/",
		controllers.UpdateScore,
	)

	router.GET( // get --> param in api
		"/courses",
		controllers.GetCourses,
	)
}
