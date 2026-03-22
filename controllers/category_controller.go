package controllers

import (
	"accounting-app/data"
	"accounting-app/models"
	"accounting-app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryData data.CategoryData
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryData: data.NewCategoryData(),
	}
}

// GetAllCategories 获取所有分类
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	categoryType := ctx.Query("type")

	categories, err := c.categoryData.GetAllCategories(categoryType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "获取分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(categories, ""))
}

// GetCategoryByID 根据ID获取分类
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	category, err := c.categoryData.GetCategoryByID(id)
	if err != nil {
		if err.Error() == "分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "获取分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(category, ""))
}

// CreateCategory 创建分类
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验必要参数
	if category.Name == "" || category.Type == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	// 校验类型
	validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
	if !validTypes[category.Type] {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的类型", "INVALID_TYPE"))
		return
	}

	created, err := c.categoryData.CreateCategory(&category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "创建分类失败", "DB_ERROR"))
		return
	}

	// 返回带子分类的结构（空）
	result := models.CategoryWithSubcategories{
		Category:      *created,
		Subcategories: []models.Subcategory{},
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(result, "创建成功"))
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验名称
	if category.Name == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少分类名称", "PARAM_ERROR"))
		return
	}

	category.ID = id

	// 校验类型（如果提供）
	if category.Type != "" {
		validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
		if !validTypes[category.Type] {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的类型", "INVALID_TYPE"))
			return
		}
	}

	updated, err := c.categoryData.UpdateCategory(&category)
	if err != nil {
		if err.Error() == "分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "更新分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(updated, "更新成功"))
}

// DeleteCategory 删除分类
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	err = c.categoryData.DeleteCategory(id)
	if err != nil {
		if err.Error() == "分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "删除分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "删除成功"))
}

// CreateSubcategory 创建子分类
func (c *CategoryController) CreateSubcategory(ctx *gin.Context) {
	var subcategory models.Subcategory
	if err := ctx.ShouldBindJSON(&subcategory); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验必要参数
	if subcategory.CategoryID == 0 || subcategory.Name == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	created, err := c.categoryData.CreateSubcategory(&subcategory)
	if err != nil {
		if err.Error() == "父分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "父分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "创建子分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(created, "创建成功"))
}

// UpdateSubcategory 更新子分类
func (c *CategoryController) UpdateSubcategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	var subcategory models.Subcategory
	if err := ctx.ShouldBindJSON(&subcategory); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验名称
	if subcategory.Name == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少子分类名称", "PARAM_ERROR"))
		return
	}

	subcategory.ID = id

	updated, err := c.categoryData.UpdateSubcategory(&subcategory)
	if err != nil {
		if err.Error() == "子分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "子分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "更新子分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(updated, "更新成功"))
}

// DeleteSubcategory 删除子分类
func (c *CategoryController) DeleteSubcategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	err = c.categoryData.DeleteSubcategory(id)
	if err != nil {
		if err.Error() == "子分类不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "子分类不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "删除子分类失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "删除成功"))
}
