package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/frontdesk/handler"
	"github.com/joshuaetim/frontdesk/infrastructure"
	"github.com/joshuaetim/frontdesk/middleware"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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
