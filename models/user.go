package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Token    string `json:"token" form:"token"`
}

func (user *User) ValidatorSanitizer() error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
