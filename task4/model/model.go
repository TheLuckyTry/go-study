package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}
