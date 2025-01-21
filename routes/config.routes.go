package routes

import (
	"net/http"
	"student_manager/controllers"
	"student_manager/middlewares"
	"student_manager/models"
	"student_manager/services"

	"github.com/gin-gonic/gin"
	//"github.com/swaggo/swag/example/basic/docs"
	//"github.com/swaggo/swag/example/basic/docs"
	//ginSwagger "github.com/swaggo/gin-swagger"
	//"github.com/swaggo/swag/example/basic/docs"
)

func ConfigRoute() *gin.Engine {
	r := gin.New()
	initRoute(r)

	// r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	// r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	// r.Use(middlewares.CORSMiddleware())
	r.POST("/admin/login", controllers.LoginAdmin)
	r.POST("/admin/logout", controllers.LogoutAdmin)
	r.POST("/admin/refreshToken", controllers.RefreshToken)

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		StudentRoute(admin)
		TeacherRoute(admin)
		ClassRoute(admin)
		// SemesterRoute(admin)
		CourseRoute(admin)
		SubjectRoute(admin)
	}

	r.POST("/teacher/login", controllers.LoginTeacher)
	r.POST("/teacher/logout", controllers.LogoutTeacher)

	teacher := r.Group("/teacher")
	teacher.Use(middlewares.TeacherMiddleware())
	{
		StudentRoute(teacher)
	}

	//docs.SwaggerInfo.BasePath = admin.BasePath() // adds /admin to swagger base path

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
	// do some other initialization staff
}
