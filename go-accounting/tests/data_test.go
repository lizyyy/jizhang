package tests

import (
	"go-accounting/data"
	"go-accounting/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFreshTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Category{}, &models.Subcategory{}, &models.Record{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestCategoryRepository_Create(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	category := &models.Category{
		Name: "测试分类",
		Type: "expense",
	}

	err := repo.Create(category)
	assert.NoError(t, err)
	assert.Greater(t, category.ID, 0)
}

func TestCategoryRepository_GetByID(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	category := &models.Category{
		Name: "测试分类",
		Type: "expense",
	}
	repo.Create(category)

	found, err := repo.GetByID(category.ID)
	assert.NoError(t, err)
	assert.Equal(t, category.Name, found.Name)
	assert.Equal(t, category.Type, found.Type)
}

func TestCategoryRepository_GetAll(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	categories := []*models.Category{
		{Name: "餐饮", Type: "expense"},
		{Name: "工资", Type: "income"},
	}
	for _, c := range categories {
		repo.Create(c)
	}

	result, err := repo.GetAll("")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = repo.GetAll("expense")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "餐饮", result[0].Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	category := &models.Category{
		Name: "测试分类",
		Type: "expense",
	}
	repo.Create(category)

	category.Name = "更新后的分类"
	err := repo.Update(category)
	assert.NoError(t, err)

	found, _ := repo.GetByID(category.ID)
	assert.Equal(t, "更新后的分类", found.Name)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	category := &models.Category{
		Name: "测试分类",
		Type: "expense",
	}
	repo.Create(category)

	err := repo.Delete(category.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(category.ID)
	assert.Error(t, err)
}

func TestSubcategory_CRUD(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewCategoryRepository(db)

	category := &models.Category{
		Name: "餐饮",
		Type: "expense",
	}
	repo.Create(category)

	subcategory := &models.Subcategory{
		CategoryID: category.ID,
		Name:       "早餐",
	}
	err := repo.CreateSubcategory(subcategory)
	assert.NoError(t, err)
	assert.Greater(t, subcategory.ID, 0)

	found, err := repo.GetSubcategoryByID(subcategory.ID)
	assert.NoError(t, err)
	assert.Equal(t, "早餐", found.Name)

	found.Name = "午餐"
	err = repo.UpdateSubcategory(found)
	assert.NoError(t, err)

	err = repo.DeleteSubcategory(found.ID)
	assert.NoError(t, err)

	_, err = repo.GetSubcategoryByID(found.ID)
	assert.Error(t, err)
}

func TestRecordRepository_Create(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewRecordRepository(db)

	category := &models.Category{
		Name: "餐饮",
		Type: "expense",
	}
	db.Create(category)

	record := &models.Record{
		Type:       "expense",
		CategoryID: category.ID,
		Amount:     100.50,
		Date:       "2024-01-15",
		Note:       "午餐",
	}

	err := repo.Create(record)
	assert.NoError(t, err)
	assert.Greater(t, record.ID, 0)
}

func TestRecordRepository_GetAll(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewRecordRepository(db)

	category := &models.Category{
		Name: "餐饮",
		Type: "expense",
	}
	db.Create(category)

	records := []*models.Record{
		{Type: "expense", CategoryID: category.ID, Amount: 100, Date: "2024-01-15"},
		{Type: "income", CategoryID: category.ID, Amount: 5000, Date: "2024-01-10"},
	}
	for _, r := range records {
		repo.Create(r)
	}

	result, err := repo.GetAll(data.RecordFilter{})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = repo.GetAll(data.RecordFilter{Type: "expense"})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "expense", result[0].Type)
}

func TestRecordRepository_GetStatistics(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewRecordRepository(db)

	category := &models.Category{
		Name: "测试",
		Type: "expense",
	}
	db.Create(category)

	records := []*models.Record{
		{Type: "expense", CategoryID: category.ID, Amount: 100, Date: "2024-01-15"},
		{Type: "expense", CategoryID: category.ID, Amount: 200, Date: "2024-01-16"},
		{Type: "income", CategoryID: category.ID, Amount: 5000, Date: "2024-01-10"},
	}
	for _, r := range records {
		repo.Create(r)
	}

	stats, err := repo.GetStatistics("", "")
	assert.NoError(t, err)
	assert.Equal(t, 300.0, stats.Expense)
	assert.Equal(t, 5000.0, stats.Income)
}

func TestRecordRepository_Delete(t *testing.T) {
	db := setupFreshTestDB(t)
	repo := data.NewRecordRepository(db)

	category := &models.Category{
		Name: "餐饮",
		Type: "expense",
	}
	db.Create(category)

	record := &models.Record{
		Type:       "expense",
		CategoryID: category.ID,
		Amount:     100,
		Date:       "2024-01-15",
	}
	repo.Create(record)

	err := repo.Delete(record.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(record.ID)
	assert.Error(t, err)
}
