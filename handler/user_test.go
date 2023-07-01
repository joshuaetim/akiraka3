package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/factory"
	"github.com/joshuaetim/frontdesk/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestUserCreateUser(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	reqBody := `{"first_name": "James", "last_name": "Harden", "email": "james@gmail.com", "password": "password", "company": "Rich Holdings"}`

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer([]byte(reqBody)))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	fmt.Println(w.Body.String())

	var result map[string]model.User
	json.NewDecoder(w.Body).Decode(&result)

	assert.EqualValues(t, "james@gmail.com", result["data"].Email)
	assert.NotEqualValues(t, "password", result["data"].Password)
	assert.EqualValues(t, 1, result["data"].ID)
}

func TestUserSignInUserWrong(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	// wrong details
	email := "josh@gmail.com"
	password := "password"
	reqBody := `{"email":"` + email + `", "password":"` + password + `"}`

	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer([]byte(reqBody)))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	fmt.Println(w.Body.String())
}

func TestUserSignInUserCorrect(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	user, err := factory.SeedUser(db)
	assert.Nil(t, err)

	// correct details
	email := user.Email
	password := "password"
	reqBody := `{"email":"` + email + `", "password":"` + password + `"}`

	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer([]byte(reqBody)))
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// fmt.Println(req.)

	fmt.Println(w.Body.String())

}

func TestUserGetUser(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	user, err := factory.SeedUser(db)
	assert.Nil(t, err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/user/%d", user.ID), nil)
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserUpdateUser(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	user, err := factory.SeedUser(db)
	assert.Nil(t, err)

	reqBody := `{"email":"` + "jjj@gmail.com" + `"}`
	req, err := http.NewRequest("PUT", fmt.Sprintf("/user/%d", user.ID), bytes.NewBuffer([]byte(reqBody)))
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}

func TestUserDeleteUser(t *testing.T) {
	initTest()
	db := infrastructure.DB()
	router := setupRoutes(db)

	w := httptest.NewRecorder()

	user, err := factory.SeedUser(db)
	assert.Nil(t, err)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/user/%d", user.ID), nil)
	assert.Nil(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
