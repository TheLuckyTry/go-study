package main

import (
	"fmt"
	"go-study/constant"
	dbbase "go-study/task3/db-base"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Students struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"column:name"`
	Age   int    `gorm:"column:age"`
	Grade string `gorm:"column:grade"`
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	db := dbbase.ConnectDB()
	err := db.AutoMigrate(&Students{})
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectDB 连接数据库
func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(constant.DBPATH), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := InitDB()
	fmt.Printf("使用的数据库文件：%s\n", constant.DBPATH)

	//编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	db.Create(&Students{Name: "张三", Age: 20, Grade: "三年级"})
	//编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	var students []Students
	db.Where("age > ?", 18).Find(&students)
	fmt.Println("所有年龄大于 18 岁的学生：", students)

	//编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	db.Model(&Students{}).Where("name = ?", "张三").Update("grade", "四年级")
	//编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	db.Where("age < ?", 15).Delete(&Students{})
	//查询所有数据
	db.Find(&students)
	fmt.Println("所有数据：", students)
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取底层数据库连接失败：" + err.Error())
	}
	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}
