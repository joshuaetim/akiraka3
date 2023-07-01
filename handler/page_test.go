package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/joshuaetim/frontdesk/factory"
	"github.com/joshuaetim/frontdesk/handler"
	"github.com/joshuaetim/frontdesk/infrastructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestPagePingServer(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, "pong", result["msg"])
}

func TestPageJsonForm(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	requestBody := []byte(`{"name":"joshua"}`)
	req, _ := http.NewRequest("POST", "/json", bytes.NewBuffer(requestBody))
	router.ServeHTTP(w, req)

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)

	assert.EqualValues(t, "joshua", result["data"])

}

func TestRequestInf(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	user, err := factory.SeedUser(db)
	assert.Nil(t, err)

	reqBody := `{"email":"` + user.Email + `"}`
	fmt.Println(reqBody)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/user/%d", user.ID), bytes.NewBuffer([]byte(reqBody)))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}

func initTest() {
	err := godotenv.Load(basepath + "/../" + ".env_test")
	fmt.Println(basepath)
	if err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", handler.Ping)
	r.POST("/json", handler.JSONRequest)

	userHandler := handler.NewUserHandler(db)
	r.POST("/user", userHandler.CreateUser)
	r.POST("/signin", userHandler.SignInUser)
	r.GET("/user/:id", userHandler.GetUser)
	r.PUT("/user/:id", userHandler.UpdateUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)

	return r
}
