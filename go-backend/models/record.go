package models

import (
	"time"
)

// Record 记账记录模型
type Record struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Type            string    `json:"type" gorm:"not null"`
	CategoryID      uint      `json:"category_id" gorm:"not null;index"`
	SubcategoryID   *uint     `json:"subcategory_id" gorm:"index"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Date            string    `json:"date" gorm:"not null"`
	Note            string    `json:"note"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// 关联字段（用于查询时展示）
	CategoryName    string `json:"category_name" gorm:"-"`
	CategoryType    string `json:"category_type" gorm:"-"`
	SubcategoryName string `json:"subcategory_name" gorm:"-"`
}

// RecordQuery 记录查询参数
type RecordQuery struct {
	Type       string
	StartDate  string
	EndDate    string
	CategoryID uint
	Keyword    string
}

// Statistics 统计结果
type Statistics struct {
	Expense  float64 `json:"expense"`
	Income   float64 `json:"income"`
	Transfer float64 `json:"transfer"`
	Loan     float64 `json:"loan"`
}

// StatisticsRow 统计查询原始结果
type StatisticsRow struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
	Count int64   `json:"count"`
}
