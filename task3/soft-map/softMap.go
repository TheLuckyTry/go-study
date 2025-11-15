package main

import (
	"fmt"
	dbbase "go-study/task3/db-base"
	"log"

	"github.com/jmoiron/sqlx"
)

type Books struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Price  int    `db:"price"`
	Author string `db:"author"`
}

func createBooksTable(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        price INTEGER NOT NULL
    );`

	_, err := db.Exec(schema)
	return err
}

// insertTestBooks 插入测试数据
func insertTestBooks(db *sqlx.DB) {
	books := []Books{
		{Title: "BOOK", Author: "作者", Price: 99},
		{Title: "BOOK2", Author: "作者2", Price: 128},
		{Title: "BOOK3", Author: "作者3", Price: 69},
		{Title: "BOOK4", Author: "作者4", Price: 89},
		{Title: "BOOK5", Author: "作者5", Price: 139},
	}

	tx := db.MustBegin()
	stmt, err := tx.Prepare("INSERT INTO books (title, author, price) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("准备插入语句失败:", err)
	}
	defer stmt.Close()

	for _, book := range books {
		_, err := stmt.Exec(book.Title, book.Author, book.Price)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			log.Fatal("插入数据失败:", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("提交事务失败:", err)
	}
}
func main() {
	db := dbbase.SqlxConnectDB()
	createBooksTable(db)
	insertTestBooks(db)
	// 查询价格大于50的书籍
	var books []Books

	// 使用参数化查询确保类型安全
	query := "SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC"
	err := db.Select(&books, query, 50)
	if err != nil {
		log.Fatal("查询书籍失败:", err)
	}

	// 输出结果
	fmt.Println("价格大于50元的书籍:")
	for _, book := range books {
		fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %d元\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
