package main

import (
	"fmt"
	"go-study/constant"
	dbbase "go-study/task3/db-base"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"column:username"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	PostCount int    `gorm:"column:post_count;default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"`
}

// Post 文章模型
type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"column:title"`
	Content       string    `gorm:"column:content"`
	UserID        uint      `gorm:"column:userId"`
	CommentCount  int       `gorm:"column:comment_count;default:0"`
	CommentStatus string    `gorm:"column:comment_status;default:'无评论'"`
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
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

// AfterCreate 在文章创建后的钩子函
func (p *Post) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("更新user:%d的文章数量\n", p.UserID)
	// 更新用户的文章数量统计
	tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count",
		tx.Select("count(*)").Model(&Post{}).Where("userId = ?", p.UserID))
	return nil
}

// AfterDelete 在评论删除后的钩子函数
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 检查文章的评论数量
	var commentCount int64
	tx.Model(&Comment{}).Where("postId = ?", c.PostID).Count(&commentCount)
	fmt.Printf("删除的评论数文章:%d的评论数量: %d\n", c.PostID, commentCount)
	// 更新文章的评论状态
	if commentCount == 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论")
		fmt.Printf("文章:%d更新为无评论", c.PostID)
	}

	// 更新文章的评论计数
	tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", commentCount)

	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// 检查文章的评论数量
	var commentCount int64
	tx.Model(&Comment{}).Where("postId = ?", c.PostID).Count(&commentCount)
	fmt.Printf("文章:%d的评论数量: %d\n", c.PostID, commentCount)
	// 更新文章的评论状态
	if commentCount > 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论")
		fmt.Printf("文章:%d更新为有评论", c.PostID)
	}

	// 更新文章的评论计数
	tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", commentCount)

	return nil
}

func main() {
	// 连接数据库
	db := dbbase.ConnectDB(constant.DBPATH)

	// 自动迁移模型，创建对应的数据库表
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		fmt.Println("创建表结构失败:", err)
	}
	fmt.Println("数据库表创建成功!")
	// 插入测试数据
	//createTestData(db)
	//查询评论数量最多的文章信息
	//queryMostCommentPost(db)
	//查询用户发布的所有文章及其对应的评论信息。
	//delCommentsByPostId(db, 1)
	//delCommentsByPostId(db, 2)
	queryAllUserContentAndPsot(db, "tom")
}
func queryAllUserContentAndPsot(db *gorm.DB, username string) User {
	var user User
	db.Transaction(func(tx *gorm.DB) error {
		tx.Preload("Posts").Preload("Posts.Comments").First(&user, "username = ?", username)
		fmt.Println("=== 用户信息 ===")
		fmt.Printf("用户名: %s, 邮箱: %s ,文章数量:%d \n", user.Username, user.Email, user.PostCount)

		fmt.Println("\n=== 文章及评论信息 ===")
		for i, post := range user.Posts {
			fmt.Printf("[%d] 文章标题: %s\n", i+1, post.Title)
			fmt.Printf("    文章内容: %s\n", post.Content)
			fmt.Printf("    文章评论数量: %d,评论状态:%s \n", post.CommentCount, post.CommentStatus)

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
		return nil
	})
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
	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{
			Username: "tom",
			Email:    "tom@qq.com",
			Password: "123456",
		}
		result := tx.FirstOrCreate(&user, User{Username: "tom"})
		if result.Error != nil {
			return result.Error
		}

		// 创建文章
		post := Post{
			Title:   "读书心得",
			Content: "这是我最近看书的读书心得",
			UserID:  user.ID,
		}
		result = tx.FirstOrCreate(&post, Post{Title: "读书心得"})
		if result.Error != nil {
			return result.Error
		}

		post2 := Post{
			Title:   "游戏攻略",
			Content: "这是一个游戏攻略",
			UserID:  user.ID,
		}
		result = tx.FirstOrCreate(&post2, Post{Title: "游戏攻略"})
		if result.Error != nil {
			return result.Error
		}

		// 创建评论
		comments := []Comment{
			{Content: "666", UserID: user.ID, PostID: post.ID},
			{Content: "博主很厉害", UserID: user.ID, PostID: post.ID},
			{Content: "666", UserID: user.ID, PostID: post2.ID},
			{Content: "游戏攻略很详细", UserID: user.ID, PostID: post2.ID},
			{Content: "是大大的", UserID: user.ID, PostID: post2.ID},
		}

		result = tx.Create(&comments)
		if result.Error != nil {
			return result.Error
		}

		fmt.Println("测试数据创建完成!")
		return nil
	})

	if err != nil {
		fmt.Printf("创建测试数据失败: %v\n", err)
	}
}

func addComments(db *gorm.DB, comments []Comment) error {
	if len(comments) == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&comments)
		if result.Error != nil {
			fmt.Printf("添加评论失败: %v\n", result.Error)
			return result.Error
		}
		return nil
	})
}

func delCommentsByPostId(db *gorm.DB, postId uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("postId = ?", postId).Delete(&Comment{})
		if result.Error != nil {
			fmt.Printf("删除文章:%d的所有评论失败: %v\n", postId, result.Error)
			return result.Error
		}

		fmt.Printf("成功删除文章:%d的%d条评论\n", postId, result.RowsAffected)
		return nil
	})
}
