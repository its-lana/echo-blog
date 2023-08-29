package controllers

import (
	"echo-blog/config"
	"echo-blog/helper"
	"echo-blog/lib/database"
	"echo-blog/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(c echo.Context) error {
	users, e := database.GetAllUsers()
	if e != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return helper.WrapResponse(http.StatusOK, "success get all user", &users).WriteToResponseBody(c.Response())

}

func GetUserByID(c echo.Context) error {
	id := c.Param("id")

	user, e := database.GetUserByID(id)

	if e != nil {
		return helper.WrapResponse(http.StatusBadRequest, "user not found", e.Error()).WriteToResponseBody(c.Response())
	}

	return helper.WrapResponse(http.StatusOK, "success get user by id", &user).WriteToResponseBody(c.Response())
}

func AddNewUser(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	if err := user.ValidatorSanitizer(); err != nil {
		return helper.WrapResponse(http.StatusBadRequest, err.Error(), &models.User{}).WriteToResponseBody(c.Response())
	}

	// Hash the user's password before saving it
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return helper.WrapResponse(http.StatusInternalServerError, "failed to hash password", err.Error()).WriteToResponseBody(c.Response())
	}
	user.Password = hashedPassword

	if err := config.DB.Save(&user).Error; err != nil {
		return helper.WrapResponse(http.StatusBadRequest, "failed to add new user", err.Error()).WriteToResponseBody(c.Response())
	}
	return helper.WrapResponse(http.StatusOK, "new user added successfully", &user).WriteToResponseBody(c.Response())
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UpdateUser(c echo.Context) error {

	idParams := c.Param("id")
	id, _ := strconv.Atoi(idParams)

	user := models.User{}
	c.Bind(&user)

	if rowsAff := config.DB.Model(&user).Where("id = ?", id).Updates(user).RowsAffected; rowsAff == 0 {
		return helper.WrapResponse(http.StatusBadRequest, "failed to update user, user id not found", &models.User{}).WriteToResponseBody(c.Response())
	}

	return helper.WrapResponse(http.StatusOK, "user updated successfully", &user).WriteToResponseBody(c.Response())
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	deletedUser, e := database.DeleteUserByID(id)

	if e != nil {
		return helper.WrapResponse(http.StatusBadRequest, "failed to delete user, id not found", e.Error()).WriteToResponseBody(c.Response())
	}
	return helper.WrapResponse(http.StatusOK, "user deleted successfully", deletedUser).WriteToResponseBody(c.Response())
}

func LoginUser(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)
	users, e := database.LoginUser(&user)

	if e != nil {
		fmt.Println(e)
		if e.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong email or password")
		} else if e.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong email or password")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}
	user.Password = ""
	return helper.WrapResponse(http.StatusOK, "login successfully", &users).WriteToResponseBody(c.Response())

}
