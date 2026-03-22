package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting-app/models"
	"accounting-app/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Response 统一响应结构（与原 Node.js 项目对齐）
type Response struct {
	Errno     int             `json:"Errno"`
	Errmsg    string          `json:"Errmsg"`
	ErrorCode string          `json:"ErrorCode"`
	Data      json.RawMessage `json:"data"`
}

func setupTestServer(t *testing.T) (*gin.Engine, *gorm.DB) {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?mode=memory"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.Category{}, &models.Subcategory{}, &models.Record{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化默认数据
	initDefaultData(t, db)

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := routes.SetupRouter(db)

	return r, db
}

func initDefaultData(t *testing.T, db *gorm.DB) {
	// 默认支出分类
	expenseCategories := []models.Category{
		{Name: "餐饮", Type: "expense"},
		{Name: "交通", Type: "expense"},
		{Name: "购物", Type: "expense"},
		{Name: "娱乐", Type: "expense"},
	}

	// 默认收入分类
	incomeCategories := []models.Category{
		{Name: "工资", Type: "income"},
		{Name: "奖金", Type: "income"},
	}

	allCategories := append(expenseCategories, incomeCategories...)

	for _, cat := range allCategories {
		if err := db.Create(&cat).Error; err != nil {
			t.Fatalf("Failed to create default category: %v", err)
		}
	}
}

// TestAPICompatibility_Categories 测试分类接口兼容性
func TestAPICompatibility_Categories(t *testing.T) {
	router, _ := setupTestServer(t)

	t.Run("GET /api/categories", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/categories", nil)
		req.Header.Set("Accept", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		// 验证响应结构
		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}
		if resp.Errmsg != "success" {
			t.Errorf("Expected Errmsg 'success', got '%s'", resp.Errmsg)
		}
		if resp.ErrorCode != "OK" {
			t.Errorf("Expected ErrorCode 'OK', got '%s'", resp.ErrorCode)
		}
		if resp.Data == nil {
			t.Error("Expected data field to be present")
		}

		// 验证 data 是数组
		var dataArray []map[string]interface{}
		if err := json.Unmarshal(resp.Data, &dataArray); err != nil {
			t.Errorf("Expected data to be an array: %v", err)
		}

		// 验证默认分类已初始化
		if len(dataArray) == 0 {
			t.Error("Expected default categories to be initialized")
		}
	})

	t.Run("POST /api/categories", func(t *testing.T) {
		body := map[string]string{
			"name": "测试分类",
			"type": "expense",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}
	})

	t.Run("GET /api/categories/:id", func(t *testing.T) {
		// 先创建一个分类
		body := map[string]string{
			"name": "测试分类2",
			"type": "income",
		}
		jsonBody, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var createResp Response
		json.Unmarshal(w.Body.Bytes(), &createResp)
		var createdCat map[string]interface{}
		json.Unmarshal(createResp.Data, &createdCat)
		catID := int(createdCat["id"].(float64))

		// 测试获取单个分类
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/categories/%d", catID), nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}
	})

	t.Run("GET /api/categories/:id - not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/categories/99999", nil)
		router.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 404 {
			t.Errorf("Expected Errno 404, got %d", resp.Errno)
		}
		if resp.ErrorCode != "NOT_FOUND" {
			t.Errorf("Expected ErrorCode 'NOT_FOUND', got '%s'", resp.ErrorCode)
		}
	})
}

// TestAPICompatibility_Records 测试记账记录接口兼容性
func TestAPICompatibility_Records(t *testing.T) {
	router, _ := setupTestServer(t)

	// 先获取一个分类ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/categories?type=expense", nil)
	router.ServeHTTP(w, req)

	var catResp Response
	json.Unmarshal(w.Body.Bytes(), &catResp)
	var categories []map[string]interface{}
	json.Unmarshal(catResp.Data, &categories)
	if len(categories) == 0 {
		t.Fatal("No categories found")
	}
	categoryID := uint(categories[0]["id"].(float64))

	t.Run("POST /api/records", func(t *testing.T) {
		body := map[string]interface{}{
			"type":        "expense",
			"category_id": categoryID,
			"amount":      100.5,
			"date":        "2024-01-15",
			"note":        "测试记录",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/records", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}

		// 验证返回的数据包含关联字段
		var record map[string]interface{}
		if err := json.Unmarshal(resp.Data, &record); err != nil {
			t.Errorf("Failed to parse record data: %v", err)
		}

		// 验证关键字段存在
		requiredFields := []string{"id", "type", "category_id", "amount", "date", "category_name"}
		for _, field := range requiredFields {
			if _, ok := record[field]; !ok {
				t.Errorf("Expected field '%s' to be present in response", field)
			}
		}
	})

	t.Run("GET /api/records", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/records", nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}

		// 验证 data 是数组
		var records []map[string]interface{}
		if err := json.Unmarshal(resp.Data, &records); err != nil {
			t.Errorf("Expected data to be an array: %v", err)
		}
	})

	t.Run("GET /api/records/statistics", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/records/statistics", nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Errno != 0 {
			t.Errorf("Expected Errno 0, got %d", resp.Errno)
		}

		// 验证统计数据结构
		var stats map[string]interface{}
		if err := json.Unmarshal(resp.Data, &stats); err != nil {
			t.Errorf("Failed to parse stats data: %v", err)
		}

		// 验证关键字段存在
		requiredFields := []string{"expense", "income", "transfer", "loan"}
		for _, field := range requiredFields {
			if _, ok := stats[field]; !ok {
				t.Errorf("Expected field '%s' to be present in stats", field)
			}
		}
	})
}

// TestAPICompatibility_ErrorCases 测试错误情况
func TestAPICompatibility_ErrorCases(t *testing.T) {
	router, _ := setupTestServer(t)

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Missing required fields", func(t *testing.T) {
		body := map[string]string{
			"name": "测试",
			// 缺少 type 字段
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid type", func(t *testing.T) {
		body := map[string]string{
			"name": "测试",
			"type": "invalid_type",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var resp Response
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.ErrorCode != "INVALID_TYPE" {
			t.Errorf("Expected ErrorCode 'INVALID_TYPE', got '%s'", resp.ErrorCode)
		}
	})

	t.Run("Not found route", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/nonexistent", nil)
		router.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var resp Response
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.ErrorCode != "NOT_FOUND" {
			t.Errorf("Expected ErrorCode 'NOT_FOUND', got '%s'", resp.ErrorCode)
		}
	})
}

// TestAPICompatibility_CORS 测试 CORS 配置
func TestAPICompatibility_CORS(t *testing.T) {
	router, _ := setupTestServer(t)

	t.Run("OPTIONS request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/categories", nil)
		router.ServeHTTP(w, req)

		if w.Code != 204 {
			t.Errorf("Expected status 204, got %d", w.Code)
		}

		// 验证 CORS 头
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("Expected Access-Control-Allow-Origin to be *")
		}
	})
}
