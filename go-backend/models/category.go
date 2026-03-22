package models

import (
	"time"
)

// Category 分类模型
type Category struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string        `json:"name" gorm:"not null"`
	Type          string        `json:"type" gorm:"not null"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	Subcategories []Subcategory `json:"subcategories" gorm:"foreignKey:CategoryID"`
}

// Subcategory 子分类模型
type Subcategory struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID uint      `json:"category_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定子分类表名
type SubcategoryModel struct{}

func (Subcategory) TableName() string {
	return "subcategories"
}
