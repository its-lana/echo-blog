package test

import (
	"echo-blog/config"
	. "echo-blog/controllers"
	"echo-blog/lib/database/seeder"
	"echo-blog/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// init function testing
func setupUserTest(t *testing.T) {
	//load env
	if err := godotenv.Load("../.env"); err != nil {
		t.Error("Error loading .env file")
	}

	//setup database
	config.InitDB()

	// clear database
	s := seeder.NewSeeder()
	fmt.Println(s)
	s.UserDelete()
	s.UserSeed()
}

func TestLoginUserSuccess(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Email:    "test1@mail.com",
		Password: "1234",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, LoginUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "login successfully", responseBody["status"])
	dataUsers := responseBody["data"].(map[string]interface{})

	assert.Equal(t, body.Email, dataUsers["email"])
	assert.Equal(t, "", dataUsers["password"])
	assert.NotNil(t, dataUsers["token"])
	assert.NotEmpty(t, dataUsers["token"])
}

func TestLoginUserWrongEmailOrPassword(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Email:    "test1@mail.com",
		Password: "1235",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	err := LoginUser(c)
	hErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, hErr.Code)
	assert.Equal(t, "wrong email or password", hErr.Message)
}

func TestGetAllUsersSuccess(t *testing.T) {
	setupUserTest(t)
	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, GetAllUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "success get all user", responseBody["status"])
}

func TestGetAllUsersFailedDBNotConnect(t *testing.T) {
	setupUserTest(t)
	db, err := config.DB.DB()
	assert.NoError(t, err)
	assert.NoError(t, db.Close())
	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	err = GetAllUser(c)
	assert.Error(t, err)
	hErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, hErr.Code)
}

func TestAddNewUserSuccess(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Username: "Budi",
		Email:    "budi@mail.com",
		Password: "12345abc",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, AddNewUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "new user added successfully", responseBody["status"])
}

func TestAddNewUserFailedWhenUserNotInputEmail(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Username: "Dodi",
		Email:    "",
		Password: "12345",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, AddNewUser(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "email is required", responseBody["status"])
}

func TestGetUserByIdSuccess(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//test
	assert.NoError(t, GetUserByID(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "success get user by id", responseBody["status"])
}

func TestGetUserByIdNotFound(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("10")

	//test
	assert.NoError(t, GetUserByID(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "user not found", responseBody["status"])
}

func TestUpdateUserByIdSuccess(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Username: "Budiman",
		Email:    "budiman@mail.com",
		Password: "12345",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//set user id
	c.Set("userId", 1)

	//test
	assert.NoError(t, UpdateUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "user updated successfully", responseBody["status"])
}

func TestUpdateUserByIdNotFound(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.User{
		Username: "Budi",
		Email:    "budi@mail.com",
		Password: "12345",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("100")

	//test
	assert.NoError(t, UpdateUser(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "failed to update user, user id not found", responseBody["status"])
}

func TestDeleteUserByIdSuccess(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//set user id
	c.Set("userId", 1)

	//test
	assert.NoError(t, DeleteUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "user deleted successfully", responseBody["status"])
}

func TestDeleteUserByIdNotFound(t *testing.T) {
	setupUserTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("10")

	//set user id
	c.Set("userId", 10)

	//test
	assert.NoError(t, DeleteUser(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "failed to delete user, id not found", responseBody["status"])
}
