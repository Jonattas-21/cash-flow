package database

import (
	"cash-flow/src/domain/dailySummary"
	"cash-flow/src/domain/transaction"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {

	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&transaction.Transaction{})
	db.AutoMigrate(&dailySummary.DailySummary{})

	return db
}
