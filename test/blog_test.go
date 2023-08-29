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
func setupBlogTest(t *testing.T) {
	//load env
	if err := godotenv.Load("../.env"); err != nil {
		t.Error("Error loading .env file")
	}

	//setup database
	config.InitDB()

	// clear database
	s := seeder.NewSeeder()
	fmt.Println(s)
	s.BlogDelete()
	s.BlogSeed()
}

func TestGetAllBlogsSuccess(t *testing.T) {
	setupBlogTest(t)
	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, GetAllBlogs(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "success get all blog", responseBody["status"])
}

func TestGetAllBlogsFailedDBNotConnect(t *testing.T) {
	setupBlogTest(t)
	db, err := config.DB.DB()
	assert.NoError(t, err)
	assert.NoError(t, db.Close())
	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	err = GetAllBlogs(c)
	assert.Error(t, err)
	hErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, hErr.Code)
}

func TestAddNewBlogsSuccess(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.Blog{
		Title: "Test Blog 3",
		Body:  "Test Body 3",
		Slug:  "slug3",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/blogs", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, AddNewBlog(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "new blog added successfully", responseBody["status"])
}

func TestAddNewBlogsFailedWhenUserNotInputAuthor(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.Blog{
		Title: "Test Blog 3",
		Body:  "",
		Slug:  "slug3",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/blogs", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//test
	assert.NoError(t, AddNewBlog(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "body is required", responseBody["status"])
}

func TestGetBlogByIdSuccess(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//test
	assert.NoError(t, GetBlogByID(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "success get blog by id", responseBody["status"])
}

func TestGetBlogByIdNotFound(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("10")

	//test
	assert.NoError(t, GetBlogByID(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "blog not found", responseBody["status"])
}

func TestUpdateBlogByIdSuccess(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.Blog{
		Title: "Test Blog Z",
		Body:  "Tester Z",
		Slug:  "slugz",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/blogs", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//test
	assert.NoError(t, UpdateBlog(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "blog updated successfully", responseBody["status"])
}

func TestUpdateBlogByIdNotFound(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//create json body
	body := models.Blog{
		Title: "Test Blog Z",
		Body:  "Tester Z",
		Slug:  "slugz",
	}

	//setup request
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/blogs", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("100")

	//test
	assert.NoError(t, UpdateBlog(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)
	assert.Equal(t, "update failed, blog id not found", responseBody["status"])
}

func TestDeleteBlogByIdSuccess(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/blogs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("1")

	//test
	assert.NoError(t, DeleteBlog(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "blog deleted successfully", responseBody["status"])
}

func TestDeleteBlogByIdNotFound(t *testing.T) {
	setupBlogTest(t)

	//setup echo context
	e := echo.New()

	//setup request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/blogs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//set params
	c.SetParamNames("id")
	c.SetParamValues("10")

	//set user id
	c.Set("userId", 10)

	//test
	assert.NoError(t, DeleteBlog(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	bodyRes, _ := io.ReadAll(rec.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(bodyRes, &responseBody)

	assert.Equal(t, "delete failed, blog id not found", responseBody["status"])
}
