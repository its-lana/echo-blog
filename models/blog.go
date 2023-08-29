package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title string `json:"title" form:"title"`
	Body  string `json:"body" form:"body"`
	Slug  string `json:"slug" form:"slug"`
}

func (blog *Blog) ValidatorSanitizer() error {
	if blog.Title == "" {
		return fmt.Errorf("title is required")
	}
	if blog.Body == "" {
		return fmt.Errorf("body is required")
	}
	if blog.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	return nil
}
