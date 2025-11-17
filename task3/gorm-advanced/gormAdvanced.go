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
	//查询用户发布的所有文章及其对应的评论信息。
	queryAllUserContentAndPsot(db, "tom")
	//查询评论数量最多的文章信息
	queryMostCommentPost(db)
}
func queryAllUserContentAndPsot(db *gorm.DB, username string) User {
	user := User{}
	db.Preload("Posts").Preload("Posts.Comments").First(&user, "username = ?", username)
	fmt.Println("=== 用户信息 ===")
	fmt.Printf("用户名: %s, 邮箱: %s\n", user.Username, user.Email)

	fmt.Println("\n=== 文章及评论信息 ===")
	for i, post := range user.Posts {
		fmt.Printf("[%d] 文章标题: %s\n", i+1, post.Title)
		fmt.Printf("    文章内容: %s\n", post.Content)

		if len(post.Comments) == 0 {
			fmt.Println("    评论: 暂无评论")
		} else {
			fmt.Println("    评论列表:")
			for j, comment := range post.Comments {
				fmt.Printf("      [%d] %s\n", j+1, comment.Content)
			}
		}
		fmt.Println()
	}
	return user
}

func queryMostCommentPost(db *gorm.DB) {
	var result struct {
		PostID uint `gorm:"column:post_id"`
		Num    int  `gorm:"column:num"`
	}
	// 添加错误处理
	err := db.Model(&Comment{}).
		Select("postId, count(*) as num").
		Group("postId").
		Order("num DESC").
		Limit(1).
		Scan(&result).Error // 注意这里添加了 .Error

	if err != nil {
		fmt.Println("查询评论统计失败:", err)
		return
	}
	// 获取对应的文章信息
	var post Post
	db.Preload("Comments").First(&post, result.PostID)

	fmt.Printf("评论数最多的文章ID: %d, 评论数: %d\n", result.PostID, result.Num)
	fmt.Printf("文章标题: %s ,文章内容: %s\n", post.Title, post.Content)
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

	post2 := Post{
		Title:   "游戏攻略",
		Content: "这是一个游戏攻略",
		UserID:  user.ID,
	}
	db.FirstOrCreate(&post2, Post{Title: "游戏攻略"})

	comments := []Comment{
		{Content: "666", UserID: user.ID, PostID: post.ID},
		{Content: "博主很厉害", UserID: user.ID, PostID: post.ID},
		{Content: "666", UserID: user.ID, PostID: post2.ID},
		{Content: "游戏攻略很详细", UserID: user.ID, PostID: post2.ID},
		{Content: "是大大的", UserID: user.ID, PostID: post2.ID},
	}
	db.Create(comments)

	fmt.Println("测试数据创建完成!")
}
