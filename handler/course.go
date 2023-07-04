package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/akiraka3/domain/model"
	"github.com/joshuaetim/akiraka3/domain/repository"
	"github.com/joshuaetim/akiraka3/infrastructure"
	"gorm.io/gorm"
)

type CourseHandler interface {
	AddCourse(*gin.Context)
	ViewCourse(*gin.Context)
	GetCoursesByStaff(*gin.Context)
	GetAllCourses(*gin.Context)
	UpdateCourse(*gin.Context)
	DeleteCourse(*gin.Context)
}

type courseHandler struct {
	repo     repository.CourseRepository
	quizRepo repository.QuizRepository
}

func NewCourseHandler(db *gorm.DB) CourseHandler {
	return &courseHandler{
		repo:     infrastructure.NewCourseRepository(db),
		quizRepo: infrastructure.NewQuizRepository(db),
	}
}

func (ch *courseHandler) AddCourse(ctx *gin.Context) {
	var course model.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	course.Staff = uint(ctx.GetFloat64("userID"))
	c, err := ch.repo.AddCourse(course)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"data": c,
	})
}

func (ch *courseHandler) ViewCourse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	course, err := ch.repo.GetCourse(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	staffId := ctx.GetFloat64("userID")
	quizzes, _ := ch.quizRepo.GetQuizzesByMap(map[string]interface{}{"staff": uint(staffId), "course": uint(course.ID)})
	quizzes = new(model.Quiz).PublicQuizArray(quizzes)
	ctx.JSON(http.StatusOK, gin.H{
		"course":  course,
		"quizzes": quizzes,
	})
}

type courseAndQuizRes struct {
	Course model.Course `json:"course"`
	Quiz   []model.Quiz `json:"quiz"`
}

func (ch *courseHandler) GetCoursesByStaff(ctx *gin.Context) {
	staffId := ctx.GetFloat64("userID")
	courses, err := ch.repo.GetAllCoursesByStaff(uint(staffId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var res []courseAndQuizRes
	for _, course := range courses {
		quiz, err := ch.quizRepo.GetQuizzesByMap(map[string]interface{}{"staff": uint(staffId), "course": uint(course.ID)})
		if err != nil {
			continue
		}
		quiz = new(model.Quiz).PublicQuizArray(quiz)
		res = append(res, courseAndQuizRes{
			Course: course,
			Quiz:   quiz,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (ch *courseHandler) GetAllCourses(ctx *gin.Context) {
	courses, err := ch.repo.GetAllCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": courses,
	})
}

func (ch *courseHandler) UpdateCourse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var course model.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	course.ID = uint(id)
	c, err := ch.repo.UpdateCourse(course)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"data": c.PublicCourse(),
	})
}

func (ch *courseHandler) DeleteCourse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	course := model.Course{}
	course.ID = uint(id)
	err := ch.repo.DeleteCourse(course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "successfully deleted",
	})
}
