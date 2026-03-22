package config

import (
	"go-accounting/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
)

var DB *gorm.DB

func InitDatabase(dbPath string) error {
	var err error
	fullPath := filepath.Join(dbPath, "accounting.db")
	DB, err = gorm.Open(sqlite.Open(fullPath), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.Category{}, &models.Subcategory{}, &models.Record{})
	if err != nil {
		return err
	}

	initDefaultData()
	return nil
}

func initDefaultData() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	defaultCategories := []models.Category{
		{Name: "餐饮", Type: "expense"},
		{Name: "交通", Type: "expense"},
		{Name: "购物", Type: "expense"},
		{Name: "娱乐", Type: "expense"},
		{Name: "医疗", Type: "expense"},
		{Name: "教育", Type: "expense"},
		{Name: "住房", Type: "expense"},
		{Name: "其他支出", Type: "expense"},
		{Name: "工资", Type: "income"},
		{Name: "奖金", Type: "income"},
		{Name: "投资", Type: "income"},
		{Name: "兼职", Type: "income"},
		{Name: "其他收入", Type: "income"},
		{Name: "银行卡转账", Type: "transfer"},
		{Name: "支付宝转账", Type: "transfer"},
		{Name: "微信转账", Type: "transfer"},
		{Name: "借款", Type: "loan"},
		{Name: "还款", Type: "loan"},
	}

	for _, cat := range defaultCategories {
		DB.Create(&cat)
	}
}
