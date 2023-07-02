package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/frontdesk/handler"
	"github.com/joshuaetim/frontdesk/infrastructure"
	"github.com/joshuaetim/frontdesk/middleware"
)

func RunAPI(address string) error {
	db := infrastructure.DB()
	userHandler := handler.NewUserHandler(db)
	quizHandler := handler.NewQuizHandler(db)
	courseHandler := handler.NewCourseHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to Akiraka v2")
	})
	apiRoutes := r.Group("/api")

	apiRoutes.GET("/checkauth", middleware.AuthorizeJWT(), handler.CheckAuth)

	userRoutes := apiRoutes.Group("/auth")
	userRoutes.POST("/register", userHandler.CreateUser)
	userRoutes.POST("/login", userHandler.SignInUser)

	userProtectedRoutes := apiRoutes.Group("/user", middleware.AuthorizeJWT())
	userProtectedRoutes.GET("/:id", userHandler.GetUser)
	userProtectedRoutes.PUT("/", userHandler.UpdateUser)

	quizRoutes := apiRoutes.Group("/quiz", middleware.AuthorizeJWT())
	quizRoutes.GET("/", quizHandler.GetAllQuiz)
	quizRoutes.POST("/course/:course", quizHandler.AddQuiz)
	quizRoutes.GET("/:id", quizHandler.ViewQuiz)
	quizRoutes.GET("/staff", quizHandler.GetQuizByStaff)
	quizRoutes.POST("/:quiz/submit", quizHandler.GradeQuiz)
	quizRoutes.DELETE("/:id", quizHandler.DeleteQuiz)

	courseRoutes := apiRoutes.Group("/courses", middleware.AuthorizeJWT())
	courseRoutes.POST("/", courseHandler.AddCourse)
	courseRoutes.GET("/", courseHandler.GetAllCourses)
	courseRoutes.GET("/:id", courseHandler.ViewCourse)
	courseRoutes.PUT("/:id", courseHandler.UpdateCourse)
	courseRoutes.DELETE("/:id", courseHandler.DeleteCourse)
	courseRoutes.GET("/staff/:staffid", courseHandler.GetCoursesByStaff)

	dashboardRoutes := apiRoutes.Group("/dashboard", middleware.AuthorizeJWT())
	dashboardRoutes.GET("/users/count", dashboardHandler.GetUsersCount)

	return r.Run(address)
}
