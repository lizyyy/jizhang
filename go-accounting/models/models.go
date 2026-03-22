package models

import "time"

type Category struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	Subcategories []Subcategory `json:"subcategories" gorm:"foreignKey:CategoryID"`
}

type Subcategory struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID int       `json:"category_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Record struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Type           string    `json:"type" gorm:"not null"`
	CategoryID     int       `json:"category_id" gorm:"not null;index"`
	SubcategoryID  *int      `json:"subcategory_id"`
	Amount         float64   `json:"amount" gorm:"not null"`
	Date           string    `json:"date" gorm:"not null"`
	Note           string    `json:"note"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type RecordWithCategory struct {
	ID             int     `json:"id"`
	Type           string  `json:"type"`
	CategoryID     int     `json:"category_id"`
	SubcategoryID  *int    `json:"subcategory_id"`
	Amount         float64 `json:"amount"`
	Date           string  `json:"date"`
	Note           string  `json:"note"`
	CreatedAt      string  `json:"created_at"`
	CategoryName   string  `json:"category_name"`
	CategoryType   string  `json:"category_type"`
	SubcategoryName string `json:"subcategory_name"`
}

type StatisticItem struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
	Count int     `json:"count"`
}

type Statistics struct {
	Expense float64 `json:"expense"`
	Income  float64 `json:"income"`
	Transfer float64 `json:"transfer"`
	Loan    float64 `json:"loan"`
}
