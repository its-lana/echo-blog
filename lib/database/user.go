package database

import (
	"echo-blog/config"
	"echo-blog/middlewares"
	"echo-blog/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() (interface{}, error) {
	var users []models.User
	if e := config.DB.Find(&users).Error; e != nil {
		return nil, e
	}
	return users, nil
}

func GetUserByID(id string) (interface{}, error) {
	var user models.User

	if e := config.DB.First(&user, id).Error; e != nil {
		return nil, e
	}
	return user, nil
}

func DeleteUserByID(id string) (interface{}, error) {
	var user models.User

	if rowsAff := config.DB.Delete(&user, id).RowsAffected; rowsAff == 0 {
		return nil, errors.New("delete failed, user id not found")
	}
	return user, nil
}

func LoginUser(user *models.User) (interface{}, error) {
	var err error
	foundUser := models.User{}

	if err = config.DB.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	user.Token, err = middlewares.CreateToken(int(foundUser.ID))
	if err != nil {
		return nil, err
	}

	if err := config.DB.Model(&foundUser).Update("token", user.Token).Error; err != nil {
		return nil, err
	}

	return user, nil
}
