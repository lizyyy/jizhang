package controllers

import (
	"go-accounting/data"
	"go-accounting/models"
	"go-accounting/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	repo *data.CategoryRepository
}

func NewCategoryController(repo *data.CategoryRepository) *CategoryController {
	return &CategoryController{repo: repo}
}

func (ctrl *CategoryController) GetAll(c *gin.Context) {
	categoryType := c.Query("type")

	categories, err := ctrl.repo.GetAll(categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "获取分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(categories))
}

func (ctrl *CategoryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	category, err := ctrl.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "分类不存在", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(category))
}

func (ctrl *CategoryController) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的类型", "INVALID_TYPE"))
		return
	}

	category := &models.Category{
		Name: req.Name,
		Type: req.Type,
	}

	if err := ctrl.repo.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "创建分类失败", "DB_ERROR"))
		return
	}

	category.Subcategories = []models.Subcategory{}
	c.JSON(http.StatusCreated, utils.SuccessResponseWithData(category, "创建成功"))
}

func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "缺少分类名称", "PARAM_ERROR"))
		return
	}

	category, err := ctrl.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "分类不存在", "NOT_FOUND"))
		return
	}

	category.Name = req.Name
	if req.Type != "" {
		validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
		if !validTypes[req.Type] {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的类型", "INVALID_TYPE"))
			return
		}
		category.Type = req.Type
	}

	if err := ctrl.repo.Update(category); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "更新分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(category, "更新成功"))
}

func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	if !ctrl.repo.Exists(id) {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "分类不存在", "NOT_FOUND"))
		return
	}

	if err := ctrl.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "删除分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(nil, "删除成功"))
}

func (ctrl *CategoryController) CreateSubcategory(c *gin.Context) {
	var req struct {
		CategoryID int    `json:"category_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	if !ctrl.repo.Exists(req.CategoryID) {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "父分类不存在", "NOT_FOUND"))
		return
	}

	subcategory := &models.Subcategory{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := ctrl.repo.CreateSubcategory(subcategory); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "创建子分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponseWithData(subcategory, "创建成功"))
}

func (ctrl *CategoryController) UpdateSubcategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "缺少子分类名称", "PARAM_ERROR"))
		return
	}

	subcategory, err := ctrl.repo.GetSubcategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "子分类不存在", "NOT_FOUND"))
		return
	}

	subcategory.Name = req.Name

	if err := ctrl.repo.UpdateSubcategory(subcategory); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "更新子分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(subcategory, "更新成功"))
}

func (ctrl *CategoryController) DeleteSubcategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	if _, err := ctrl.repo.GetSubcategoryByID(id); err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "子分类不存在", "NOT_FOUND"))
		return
	}

	if err := ctrl.repo.DeleteSubcategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "删除子分类失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(nil, "删除成功"))
}
