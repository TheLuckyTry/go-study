package main

import (
	"fmt"
	dbbase "go-study/task3/db-base"
	"log"

	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

// 创建 employees 表
func createEmployeesTable(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS employees (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        department TEXT NOT NULL,
        salary INTEGER NOT NULL
    );`

	_, err := db.Exec(schema)
	return err
}

// 插入测试数据
func insertTestEmployees(db *sqlx.DB) error {
	// 批量插入测试数据
	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 8000},
		{Name: "李四", Department: "技术部", Salary: 9500},
		{Name: "王五", Department: "销售部", Salary: 6000},
		{Name: "赵六", Department: "人事部", Salary: 7000},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, emp := range employees {
		_, err := stmt.Exec(emp.Name, emp.Department, emp.Salary)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func main() {
	db := dbbase.SqlxConnectDB()
	createEmployeesTable(db)
	insertTestEmployees(db)
	defer db.Close()
	// 查询部门为"技术部"的所有员工
	var employees []Employee
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatal("查询技术部员工失败:", err)
	}
	for _, emp := range employees {
		fmt.Printf("技术部员工ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	var employee Employee
	err = db.Get(&employee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		log.Fatal("查询最高工资员工失败:", err)
	}
	fmt.Printf("最高工资员工ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", employee.ID, employee.Name, employee.Department, employee.Salary)
}
