package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/domain/repository"
	"github.com/joshuaetim/frontdesk/infrastructure"
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
	repo repository.CourseRepository
}

func NewCourseHandler(db *gorm.DB) CourseHandler {
	return &courseHandler{
		repo: infrastructure.NewCourseRepository(db),
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
	ctx.JSON(http.StatusOK, gin.H{
		"data": course,
	})
}

func (ch *courseHandler) GetCoursesByStaff(ctx *gin.Context) {
	staffId, _ := strconv.Atoi(ctx.Param("staffid"))
	courses, err := ch.repo.GetAllCoursesByStaff(uint(staffId))
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
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"data": "successfully deleted",
	})
}
