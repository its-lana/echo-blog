package database

import (
	"echo-blog/config"
	"echo-blog/models"
	"errors"
)

func GetAllBlogs() (interface{}, error) {
	var blogs []models.Blog
	if e := config.DB.Find(&blogs).Error; e != nil {
		return nil, e
	}
	return blogs, nil
}

func GetBlogByID(id string) (interface{}, error) {
	var blog models.Blog

	if e := config.DB.First(&blog, id).Error; e != nil {
		return nil, e
	}
	return blog, nil
}

func DeleteBlogByID(id string) (interface{}, error) {
	var blog models.Blog

	if rowsAff := config.DB.Delete(&blog, id).RowsAffected; rowsAff == 0 {
		return nil, errors.New("delete failed, blog id not found")
	}
	return blog, nil
}
