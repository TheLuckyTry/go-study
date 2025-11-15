package main

import (
	"fmt"
	_ "fmt"
	_ "go-study/constant"
	dbbase "go-study/task3/db-base"

	_ "github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type Accounts struct {
	ID      uint            `gorm:"primaryKey"`
	Balance decimal.Decimal `gorm:"column:balance;type:decimal(10,2)"`
}

type Transactions struct {
	ID            uint            `gorm:"primaryKey"`
	FromAccountID uint            `gorm:"column:from_account_id"`
	ToAccountID   uint            `gorm:"column:to_account_id"`
	Amount        decimal.Decimal `gorm:"column:amount;type:decimal(10,2)"`
}

func transferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount decimal.Decimal) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 检查转出账户是否存在且余额足够
		var fromAccount Accounts
		result := tx.Where("id = ?", fromAccountID).First(&fromAccount)
		if result.Error != nil {
			return fmt.Errorf("查询转出账户失败: %w", result.Error)
		}

		if fromAccount.Balance.LessThan(amount) {
			return fmt.Errorf("账户余额不足，当前余额: %s", fromAccount.Balance.String())
		}

		// 检查转入账户是否存在
		var toAccount Accounts
		result = tx.Where("id = ?", toAccountID).First(&toAccount)
		if result.Error != nil {
			return fmt.Errorf("查询转入账户失败: %w", result.Error)
		}

		// 从转出账户扣除金额
		newFromBalance := fromAccount.Balance.Sub(amount)
		result = tx.Model(&fromAccount).Update("balance", newFromBalance)
		if result.Error != nil {
			return fmt.Errorf("更新转出账户失败: %w", result.Error)
		}

		// 向转入账户增加金额
		newToBalance := toAccount.Balance.Add(amount)
		result = tx.Model(&toAccount).Update("balance", newToBalance)
		if result.Error != nil {
			return fmt.Errorf("更新转入账户失败: %w", result.Error)
		}

		// 记录转账交易
		transaction := Transactions{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
		}
		result = tx.Create(&transaction)
		if result.Error != nil {
			return fmt.Errorf("创建交易记录失败: %w", result.Error)
		}

		return nil
	})
}
func InitDB() *gorm.DB {
	db := dbbase.ConnectDB()
	err := db.AutoMigrate(&Accounts{}, &Transactions{})
	if err != nil {
		panic(err)
	}
	return db
}

/*
*
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。,
,
,
*/

func main() {

	db := InitDB()
	accountA := Accounts{
		ID:      1,
		Balance: decimal.NewFromInt(200),
	}
	accountB := Accounts{
		ID:      2,
		Balance: decimal.NewFromInt(100),
	}
	db.FirstOrCreate(&accountA)
	db.FirstOrCreate(&accountB)
	showAccountBalances(db)

	// 执行转账操作
	err := transferMoney(db, 2, 1, decimal.NewFromInt(50))
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功!")
	}
	showAccountBalances(db)
	showTransactions(db)

}
func showAccountBalances(db *gorm.DB) {
	var accounts []Accounts
	db.Find(&accounts)
	fmt.Println("账户余额:")
	for _, account := range accounts {
		fmt.Printf("账户ID: %d, 余额: %s\n", account.ID, account.Balance.String())
	}
}
func showTransactions(db *gorm.DB) {
	var transactions []Transactions
	db.Find(&transactions)
	fmt.Println("交易记录:")
	for _, transaction := range transactions {
		fmt.Printf("ID: %d, FromAccountID: %d, ToAccountID: %d, Amount: %s\n",
			transaction.ID, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount.String())
	}
}
