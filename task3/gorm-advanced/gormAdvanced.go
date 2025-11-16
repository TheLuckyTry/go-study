package main

import (
	"fmt"
	dbbase "go-study/task3/db-base"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}

// Post 文章模型
type Post struct {
	ID       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"column:title"`
	Content  string    `gorm:"column:content"`
	UserID   uint      `gorm:"column:userId"`
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment 评论模型
type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"column:content"`
	UserID  uint   `gorm:"column:userId"`
	User    User   `gorm:"foreignKey:UserID"`
	PostID  uint   `gorm:"column:postId"`
	Post    Post   `gorm:"foreignKey:PostID"`
}

func main() {
	// 连接数据库
	db := dbbase.ConnectDB()

	// 自动迁移模型，创建对应的数据库表
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		fmt.Println("创建表结构失败:", err)
	}

	fmt.Println("数据库表创建成功!")

	// 插入测试数据
	createTestData(db)
}

// createTestData 创建测试数据
func createTestData(db *gorm.DB) {
	// 创建用户
	user := User{
		Username: "tom",
		Email:    "tom@qq.com",
		Password: "123456",
	}
	db.FirstOrCreate(&user, User{Username: "tom"})

	// 创建文章
	post := Post{
		Title:   "读书心得",
		Content: "这是我最近看书的读书心得",
		UserID:  user.ID,
	}
	db.FirstOrCreate(&post, Post{Title: "读书心得"})

	// 创建评论
	comment := Comment{
		Content: "666",
		UserID:  user.ID,
		PostID:  post.ID,
	}
	db.Create(&comment)
	// 创建评论
	comment2 := Comment{
		Content: "博主很厉害",
		UserID:  user.ID,
		PostID:  post.ID,
	}
	db.Create(&comment2)

	fmt.Println("测试数据创建完成!")
}
