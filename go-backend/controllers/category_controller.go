package controllers

import (
	"accounting-app/data"
	"accounting-app/models"
	"accounting-app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryData data.CategoryData
}

// NewCategoryController 创建分类控制器
func NewCategoryController(categoryData data.CategoryData) *CategoryController {
	return &CategoryController{categoryData: categoryData}
}

// GetAllCategories 获取所有分类
func (cc *CategoryController) GetAllCategories(c *gin.Context) {
	categoryType := c.Query("type")

	categories, err := cc.categoryData.GetAllCategories(categoryType)
	if err != nil {
		utils.JSONDBError(c, "获取分类失败")
		return
	}

	utils.JSONSuccess(c, categories)
}

// GetCategoryByID 获取单个分类
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的分类ID")
		return
	}

	category, err := cc.categoryData.GetCategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取分类失败")
		return
	}

	if category == nil {
		utils.JSONNotFound(c, "分类不存在")
		return
	}

	utils.JSONSuccess(c, category)
}

// CreateCategory 创建分类
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, "缺少必要参数")
		return
	}

	// 验证类型
	validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
	if !validTypes[req.Type] {
		utils.JSONError(c, 400, "无效的类型", "INVALID_TYPE")
		return
	}

	category := &models.Category{
		Name: req.Name,
		Type: req.Type,
	}

	if err := cc.categoryData.CreateCategory(category); err != nil {
		utils.JSONDBError(c, "创建分类失败")
		return
	}

	// 返回新创建的分类（包含空的子分类数组）
	category.Subcategories = []models.Subcategory{}
	utils.JSONCreated(c, category, "创建成功")
}

// UpdateCategory 更新分类
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的分类ID")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, "缺少分类名称")
		return
	}

	// 检查分类是否存在
	existingCategory, err := cc.categoryData.GetCategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取分类失败")
		return
	}
	if existingCategory == nil {
		utils.JSONNotFound(c, "分类不存在")
		return
	}

	// 验证类型（如果提供了）
	categoryType := existingCategory.Type
	if req.Type != "" {
		validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
		if !validTypes[req.Type] {
			utils.JSONError(c, 400, "无效的类型", "INVALID_TYPE")
			return
		}
		categoryType = req.Type
	}

	category := &models.Category{
		ID:   uint(id),
		Name: req.Name,
		Type: categoryType,
	}

	if err := cc.categoryData.UpdateCategory(category); err != nil {
		utils.JSONDBError(c, "更新分类失败")
		return
	}

	// 获取更新后的分类
	updatedCategory, err := cc.categoryData.GetCategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取分类失败")
		return
	}

	utils.JSONSuccess(c, updatedCategory, "更新成功")
}

// DeleteCategory 删除分类
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的分类ID")
		return
	}

	// 检查分类是否存在
	existingCategory, err := cc.categoryData.GetCategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取分类失败")
		return
	}
	if existingCategory == nil {
		utils.JSONNotFound(c, "分类不存在")
		return
	}

	if err := cc.categoryData.DeleteCategory(uint(id)); err != nil {
		utils.JSONDBError(c, "删除分类失败")
		return
	}

	utils.JSONSuccess(c, nil, "删除成功")
}

// CreateSubcategory 创建子分类
func (cc *CategoryController) CreateSubcategory(c *gin.Context) {
	var req struct {
		CategoryID uint   `json:"category_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, "缺少必要参数")
		return
	}

	// 检查父分类是否存在
	parentCategory, err := cc.categoryData.GetCategoryByID(req.CategoryID)
	if err != nil {
		utils.JSONDBError(c, "检查分类失败")
		return
	}
	if parentCategory == nil {
		utils.JSONNotFound(c, "父分类不存在")
		return
	}

	subcategory := &models.Subcategory{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := cc.categoryData.CreateSubcategory(subcategory); err != nil {
		utils.JSONDBError(c, "创建子分类失败")
		return
	}

	utils.JSONCreated(c, subcategory, "创建成功")
}

// UpdateSubcategory 更新子分类
func (cc *CategoryController) UpdateSubcategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的子分类ID")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, "缺少子分类名称")
		return
	}

	// 检查子分类是否存在
	existingSubcategory, err := cc.categoryData.GetSubcategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取子分类失败")
		return
	}
	if existingSubcategory == nil {
		utils.JSONNotFound(c, "子分类不存在")
		return
	}

	subcategory := &models.Subcategory{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := cc.categoryData.UpdateSubcategory(subcategory); err != nil {
		utils.JSONDBError(c, "更新子分类失败")
		return
	}

	// 获取更新后的子分类
	updatedSubcategory, err := cc.categoryData.GetSubcategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取子分类失败")
		return
	}

	utils.JSONSuccess(c, updatedSubcategory, "更新成功")
}

// DeleteSubcategory 删除子分类
func (cc *CategoryController) DeleteSubcategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的子分类ID")
		return
	}

	// 检查子分类是否存在
	existingSubcategory, err := cc.categoryData.GetSubcategoryByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取子分类失败")
		return
	}
	if existingSubcategory == nil {
		utils.JSONNotFound(c, "子分类不存在")
		return
	}

	if err := cc.categoryData.DeleteSubcategory(uint(id)); err != nil {
		utils.JSONDBError(c, "删除子分类失败")
		return
	}

	utils.JSONSuccess(c, nil, "删除成功")
}
