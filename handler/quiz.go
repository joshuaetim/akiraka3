package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
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
	UpdateQuiz(*gin.Context)
}

type quizHandler struct {
	repo      repository.QuizRepository
	gradeRepo repository.GradeRepository
	userRepo  repository.UserRepository
}

func NewQuizHandler(db *gorm.DB) QuizHandler {
	return &quizHandler{
		repo:      infrastructure.NewQuizRepository(db),
		gradeRepo: infrastructure.NewGradeRepository(db),
		userRepo:  infrastructure.NewUserRepository(db),
	}
}

type gradeRes struct {
	Grade model.Grade `json:"grade"`
	User  model.User  `json:"user"`
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
	// check staff
	userId := ctx.GetFloat64("userID")
	user, err := qh.userRepo.GetUser(uint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var gradeResArray []gradeRes
	if user.Role == "teacher" {
		grades, err := qh.gradeRepo.GetGradesByMap(map[string]interface{}{"quiz": uint(id)})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		for _, grade := range grades {
			user, err := qh.userRepo.GetUser(grade.User)
			if err != nil {
				continue
			}
			res := gradeRes{Grade: grade, User: *user.PublicUser()}
			gradeResArray = append(gradeResArray, res)
		}
	}
	var doneFlag bool
	if user.Role == "student" {
		grade, _ := qh.gradeRepo.GetGradesByMap(map[string]interface{}{"user": user.ID, "quiz": quiz.ID})
		if len(grade) > 0 {
			doneFlag = true
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"quiz":   quiz.PublicQuiz(),
		"grades": gradeResArray,
		"done":   doneFlag,
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

	title := ctx.Request.FormValue("title")
	duration, _ := strconv.Atoi(ctx.Request.FormValue("duration"))
	if duration == 0 {
		return
	}

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
		Title:     title,
		Duration:  duration,
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
	publicQuizzes := new(model.Quiz).PublicQuizArray(quizzes)
	ctx.JSON(http.StatusOK, gin.H{
		"data": publicQuizzes,
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
	quizzes, err := qh.repo.GetQuizzesByMap(map[string]interface{}{"staff": uint(uID)})
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
	scorePercent := (float64(score) / float64(len(actualAnswers))) * 100
	scorePercent = math.Round(scorePercent*100) / 100
	grade := model.Grade{
		Quiz:  quiz.ID,
		Score: scorePercent,
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

func (qh *quizHandler) UpdateQuiz(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	quiz := model.Quiz{}
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	quiz.ID = uint(id)

	_, err := qh.repo.UpdateQuiz(quiz)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "updated successfully",
	})
}
