package seeder

import (
	"echo-blog/config"
	"echo-blog/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type seed struct {
	DB *gorm.DB
}

func NewSeeder() *seed {
	config.InitDB()
	return &seed{DB: config.DB}
}

func (s *seed) UserSeed() {
	users := []models.User{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Username: "test1",
			Email:    "test1@mail.com",
			Password: "1234", // Original password
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Username: "test2",
			Email:    "test2@mail.com",
			Password: "1234", // Original password
		},
	}

	// Hash passwords before inserting
	for i := range users {
		hashedPassword, err := hashPassword(users[i].Password)
		if err != nil {
			log.Printf("cannot hash password for user ID %d, error : %v\n", users[i].ID, err)
			return
		}
		users[i].Password = hashedPassword
	}

	if err := s.DB.Create(&users).Error; err != nil {
		log.Printf("cannot seed data users, error : %v\n", err)
		return
	}
	log.Println("success seed data users")
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *seed) BlogSeed() {
	blogs := []models.Blog{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Title: "Test Blog 1",
			Body:  "Test Body 1",
			Slug:  "slug1",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Title: "Test Blog 2",
			Body:  "Test Body 2",
			Slug:  "slug2",
		},
	}
	if err := s.DB.Create(&blogs).Error; err != nil {
		log.Printf("cannot seed data blogs, error : %v\n", err)
	}
	log.Println("success seed data blogs")
}

func (s *seed) UserDelete() {
	s.DB.Exec("DELETE FROM users")
}

func (s *seed) BlogDelete() {
	s.DB.Exec("DELETE FROM blogs")
}
