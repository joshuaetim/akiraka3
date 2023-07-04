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

type GradeHandler interface {
	GetGrades(*gin.Context)
	GetGradesByQuiz(*gin.Context)
	GetGradesByUser(*gin.Context)
	GetGrade(*gin.Context)
	DeleteGrade(*gin.Context)
}

type gradeHandler struct {
	repo     repository.GradeRepository
	quizRepo repository.QuizRepository
}

func NewGradeHandler(db *gorm.DB) GradeHandler {
	return &gradeHandler{
		repo:     infrastructure.NewGradeRepository(db),
		quizRepo: infrastructure.NewQuizRepository(db),
	}
}

func (gh *gradeHandler) GetGrades(ctx *gin.Context) {
	grades, err := gh.repo.GetAllGrades()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": grades,
	})
}

func (gh *gradeHandler) GetGradesByQuiz(ctx *gin.Context) {
	quizId, _ := strconv.Atoi(ctx.Param("quiz"))
	grades, err := gh.repo.GetAllByQuiz(uint(quizId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"grades": grades,
	})
}

func (gh *gradeHandler) GetGrade(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	grade, err := gh.repo.GetGrade(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": grade,
	})
}

func (gh *gradeHandler) DeleteGrade(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	grade := model.Grade{}
	grade.ID = uint(id)
	err := gh.repo.DeleteGrade(grade)
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

type gradeWithQuizInfo struct {
	Grade model.Grade `json:"grade"`
	Quiz  model.Quiz  `json:"quiz"`
}

func (gh *gradeHandler) GetGradesByUser(ctx *gin.Context) {
	userId := ctx.GetFloat64("userID")
	grades, err := gh.repo.GetByUser(uint(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var res []gradeWithQuizInfo
	for _, grade := range grades {
		quiz, err := gh.quizRepo.GetQuiz(grade.Quiz)
		if err != nil {
			continue
		}
		res = append(res, gradeWithQuizInfo{
			Grade: grade,
			Quiz: model.Quiz{
				Title: quiz.Title,
			},
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
