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

// RequestLog 请求日志表
type RequestLog struct {
	gorm.Model
	Method     string `gorm:"size:10;not null"`  // 请求方法
	Path       string `gorm:"size:255;not null"` // 请求路径
	StatusCode int    `gorm:"not null"`          // 状态码
	UserID     uint   `gorm:"index"`             // 用户ID（如果有）
	Username   string `gorm:"size:100;index"`    // 用户名（如果有）
	IPAddress  string `gorm:"size:45"`           // IP地址
	UserAgent  string `gorm:"size:500"`          // 用户代理
	Request    string `gorm:"type:text"`         // 请求参数
	Response   string `gorm:"type:text"`         // 响应数据
	Error      string `gorm:"type:text"`         // 错误信息
	Duration   int64  `gorm:"not null"`          // 耗时（毫秒）
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{}, &RequestLog{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}
