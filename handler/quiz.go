package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/domain/repository"
	"github.com/joshuaetim/frontdesk/infrastructure"
	"gorm.io/gorm"
)

type QuizHandler interface {
	ViewQuiz(*gin.Context)
	SubmitQuiz(*gin.Context)
	AddQuiz(*gin.Context)
	GetAllQuiz(*gin.Context)
	DeleteQuiz(*gin.Context)
	// GetQuiz(*gin.Context)
	GetQuizByStaff(*gin.Context)
	GradeQuiz(*gin.Context)
}

type quizHandler struct {
	repo      repository.QuizRepository
	gradeRepo repository.GradeRepository
}

func NewQuizHandler(db *gorm.DB) QuizHandler {
	return &quizHandler{
		repo:      infrastructure.NewQuizRepository(db),
		gradeRepo: infrastructure.NewGradeRepository(db),
	}
}

func (qh *quizHandler) ViewQuiz(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	quiz, err := qh.repo.GetQuiz(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": quiz,
	})
}

type quizAnswer struct {
	Answers []string `json:"answers"`
}

func (qh *quizHandler) SubmitQuiz(ctx *gin.Context) {
	// userId := ctx.GetFloat64("userID")
	var qA quizAnswer
	if err := ctx.ShouldBindJSON(&qA); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "there was an error in your request",
		})
		return
	}

	// quizId, _ := strconv.Atoi(ctx.Param("quiz"))
	// quiz, err := uh.quizRepo
}

var Folder = "/quizzes"

func UploadFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fullpath := fmt.Sprintf("%s/%d%s", Folder, time.Now().Unix(), filepath.Base(header.Filename))
	newFile, err := os.Create(fullpath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	io.Copy(newFile, file)
}

const SEP = "===SEP==="

func (qh *quizHandler) AddQuiz(ctx *gin.Context) {
	// get csv data
	file, _, _ := ctx.Request.FormFile("file")
	csvReader := csv.NewReader(file)

	var questions []string
	var answers []string
	var options []string

	for i := 0; ; i++ {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if i == 0 {
			continue
		}

		for i, field := range rec {
			if i == 0 {
				questions = append(questions, field)
			}
			if i == 1 {
				options = append(options, field)
			}
			if i == 2 {
				answers = append(answers, field)
			}
		}
	}
	questionsCSV := strings.Join(questions, SEP)
	answersCSV := strings.Join(answers, SEP)
	optionsCSV := strings.Join(options, SEP)
	courseId, _ := strconv.Atoi(ctx.Param("course"))

	staffId := ctx.GetFloat64("userID")

	quiz := model.Quiz{
		Questions: questionsCSV,
		Answers:   answersCSV,
		Options:   optionsCSV,
		Course:    uint(courseId),
		Staff:     uint(staffId),
	}
	_, err := qh.repo.AddQuiz(quiz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "problem parsing your request",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"data": "quiz created successful",
	})

}

func (qh *quizHandler) GetAllQuiz(ctx *gin.Context) {
	quizzes, err := qh.repo.GetAllQuizzes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": quizzes,
	})
}

func (qh *quizHandler) DeleteQuiz(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	quiz := model.Quiz{}
	quiz.ID = uint(id)

	err := qh.repo.DeleteQuiz(quiz)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "quiz deleted success",
	})
}

func (qh *quizHandler) GetQuizByStaff(ctx *gin.Context) {
	uID := ctx.GetFloat64("userID")
	quizzes, err := qh.repo.GetQuizzesByStaff(uint(uID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": quizzes,
	})
}

type quizRequest struct {
	Answers []string `json:"answers"`
}

func (qh *quizHandler) GradeQuiz(ctx *gin.Context) {
	userId := ctx.GetFloat64("userID")
	var quizReq quizRequest
	if err := ctx.ShouldBindJSON(&quizReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	quizId, _ := strconv.Atoi(ctx.Param("quiz"))
	// check if submitted
	if _, err := qh.gradeRepo.GetByUserAndQuiz(uint(userId), uint(quizId)); err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "you have submitted already",
		})
		return
	}
	quiz, err := qh.repo.GetQuiz(uint(quizId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var score int
	actualAnswers := strings.Split(quiz.Answers, SEP)
	for i, ans := range quizReq.Answers {
		if i >= len(actualAnswers) {
			break
		}
		if ans == actualAnswers[i] {
			score++
		}
	}
	grade := model.Grade{
		Quiz:  quiz.ID,
		Score: score,
		User:  uint(userId),
	}
	err = qh.gradeRepo.AddGrade(grade)
	if err != nil {
		msg := err.Error()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "graded successfully",
	})
}
