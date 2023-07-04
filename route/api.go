package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/akiraka3/handler"
	"github.com/joshuaetim/akiraka3/infrastructure"
	"github.com/joshuaetim/akiraka3/middleware"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RunAPI(address string) error {
	db := infrastructure.DB()
	userHandler := handler.NewUserHandler(db)
	quizHandler := handler.NewQuizHandler(db)
	courseHandler := handler.NewCourseHandler(db)
	gradeHandler := handler.NewGradeHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to Akiraka v2")
	})
	apiRoutes := r.Group("/api")

	apiRoutes.GET("/checkauth", middleware.AuthorizeJWT(), handler.CheckAuth)

	userRoutes := apiRoutes.Group("/auth")
	userRoutes.POST("/register", userHandler.CreateUser)
	userRoutes.POST("/login", userHandler.SignInUser)

	apiRoutes.GET("/users", middleware.AuthorizeJWT(), userHandler.GetUsers)

	userProtectedRoutes := apiRoutes.Group("/user", middleware.AuthorizeJWT())
	userProtectedRoutes.GET("/:id", userHandler.GetUser)
	userProtectedRoutes.PUT("/", userHandler.UpdateUser)
	userProtectedRoutes.GET("/", userHandler.GetCurrentUser)

	quizRoutes := apiRoutes.Group("/quiz", middleware.AuthorizeJWT())
	quizRoutes.GET("/", quizHandler.GetAllQuiz)
	quizRoutes.POST("/course/:course", quizHandler.AddQuiz)
	quizRoutes.GET("/:id", quizHandler.ViewQuiz)
	quizRoutes.GET("/staff", quizHandler.GetQuizByStaff)
	quizRoutes.POST("/:quiz/submit", quizHandler.GradeQuiz)
	quizRoutes.PATCH("/:id", quizHandler.UpdateQuiz)
	quizRoutes.DELETE("/:id", quizHandler.DeleteQuiz)

	gradeRoutes := apiRoutes.Group("/grade", middleware.AuthorizeJWT())
	gradeRoutes.GET("/", gradeHandler.GetGrades)
	gradeRoutes.GET("/:id", gradeHandler.GetGrade)
	gradeRoutes.GET("/quiz/:quiz", gradeHandler.GetGradesByQuiz)
	gradeRoutes.GET("/user/", gradeHandler.GetGradesByUser)
	gradeRoutes.DELETE("/:id", gradeHandler.DeleteGrade)

	courseRoutes := apiRoutes.Group("/courses", middleware.AuthorizeJWT())
	courseRoutes.POST("/", courseHandler.AddCourse)
	courseRoutes.GET("/", courseHandler.GetAllCourses)
	courseRoutes.GET("/:id", courseHandler.ViewCourse)
	courseRoutes.PUT("/:id", courseHandler.UpdateCourse)
	courseRoutes.DELETE("/:id", courseHandler.DeleteCourse)
	courseRoutes.GET("/staff", courseHandler.GetCoursesByStaff)

	dashboardRoutes := apiRoutes.Group("/dashboard", middleware.AuthorizeJWT())
	dashboardRoutes.GET("/users/count", dashboardHandler.GetUsersCount)

	return r.Run(address)
}
