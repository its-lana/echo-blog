package controllers

import (
	"echo-blog/config"
	"echo-blog/helper"
	"echo-blog/lib/database"
	"echo-blog/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllBlogs(c echo.Context) error {
	blogs, e := database.GetAllBlogs()
	if e != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return helper.WrapResponse(http.StatusOK, "success get all blog", &blogs).WriteToResponseBody(c.Response())

}

func GetBlogByID(c echo.Context) error {
	id := c.Param("id")

	blog, e := database.GetBlogByID(id)

	if e != nil {
		return helper.WrapResponse(http.StatusBadRequest, "blog not found", e.Error()).WriteToResponseBody(c.Response())
	}

	return helper.WrapResponse(http.StatusOK, "success get blog by id", &blog).WriteToResponseBody(c.Response())
}

func AddNewBlog(c echo.Context) error {
	blog := models.Blog{}
	c.Bind(&blog)

	if err := blog.ValidatorSanitizer(); err != nil {
		return helper.WrapResponse(http.StatusBadRequest, err.Error(), &models.Blog{}).WriteToResponseBody(c.Response())
	}

	if err := config.DB.Save(&blog).Error; err != nil {
		return helper.WrapResponse(http.StatusBadRequest, "failed to add new blog", err.Error()).WriteToResponseBody(c.Response())
	}
	return helper.WrapResponse(http.StatusOK, "new blog added successfully", &blog).WriteToResponseBody(c.Response())
}

func UpdateBlog(c echo.Context) error {

	id := c.Param("id")

	blog := models.Blog{}
	c.Bind(&blog)

	if rowsAff := config.DB.Model(&blog).Where("id = ?", id).Updates(blog).RowsAffected; rowsAff == 0 {
		return helper.WrapResponse(http.StatusBadRequest, "update failed, blog id not found", &models.Blog{}).WriteToResponseBody(c.Response())
	}

	return helper.WrapResponse(http.StatusOK, "blog updated successfully", &blog).WriteToResponseBody(c.Response())
}

func DeleteBlog(c echo.Context) error {
	id := c.Param("id")

	_, e := database.DeleteBlogByID(id)

	if e != nil {
		return helper.WrapResponse(http.StatusBadRequest, "delete failed, blog id not found", e.Error()).WriteToResponseBody(c.Response())
	}
	return helper.WrapResponse(http.StatusOK, "blog deleted successfully", &models.Blog{}).WriteToResponseBody(c.Response())
}
