package models

import (
	"time"
)

// Category 分类模型
type Category struct {
	ID           int       `json:"id" gorm:"primaryKey;column:id"`
	Name         string    `json:"name" gorm:"column:name"`
	Type         string    `json:"type" gorm:"column:type"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
}

// Subcategory 子分类模型
type Subcategory struct {
	ID         int       `json:"id" gorm:"primaryKey;column:id"`
	CategoryID int       `json:"category_id" gorm:"column:category_id"`
	Name       string    `json:"name" gorm:"column:name"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

// Record 记账记录模型
type Record struct {
	ID            int       `json:"id" gorm:"primaryKey;column:id"`
	Type          string    `json:"type" gorm:"column:type"`
	CategoryID    int       `json:"category_id" gorm:"column:category_id"`
	SubcategoryID *int      `json:"subcategory_id" gorm:"column:subcategory_id"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	Date          string    `json:"date" gorm:"column:date"`
	Note          string    `json:"note" gorm:"note`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
}

// RecordWithDetails 包含关联信息的记录
type RecordWithDetails struct {
	Record
	CategoryName    string `json:"category_name"`
	CategoryType  string `json:"category_type"`
	SubcategoryName *string `json:"subcategory_name"`
}

// CategoryWithSubcategories 带子分类的分类
type CategoryWithSubcategories struct {
	Category
	Subcategories []Subcategory `json:"subcategories"`
}

// Statistics 统计数据
type Statistics struct {
	Expense float64 `json:"expense"`
	Income  float64 `json:"income"`
	Transfer float64 `json:"transfer"`
	Loan     float64 `json:"loan"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

func (Subcategory) TableName() string {
	return "subcategories"
}

func (Record) TableName() string {
	return "records"
}
