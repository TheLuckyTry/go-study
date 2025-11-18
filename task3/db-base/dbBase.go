package dbbase

import (
	"go-study/constant"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

func ConnectDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func SqlxConnectDB() *sqlx.DB {
	// 连接数据库
	db, err := sqlx.Connect("sqlite", constant.DBPATH)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	return db
}
