package main

import (
	_ "fmt"
	_ "go-study/constant"

	_ "github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	_ "gorm.io/gorm"
)

type Accounts struct {
	ID      uint            `gorm:"primaryKey"`
	Balance decimal.Decimal `gorm:"column:balance;type:decimal(10,2)"`
}

type Transactions struct {
	ID            uint            `gorm:"primaryKey"`
	fromAccountId uint            `gorm:"column:from_account_id"`
	toAccountId   uint            `gorm:"column:to_account_id"`
	amount        decimal.Decimal `gorm:"column:amount;type:decimal(10,2)"`
}

func main() {

}
