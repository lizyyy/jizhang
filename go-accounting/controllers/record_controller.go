package controllers

import (
	"go-accounting/data"
	"go-accounting/models"
	"go-accounting/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecordController struct {
	repo *data.RecordRepository
}

func NewRecordController(repo *data.RecordRepository) *RecordController {
	return &RecordController{repo: repo}
}

func (ctrl *RecordController) GetAll(c *gin.Context) {
	filter := data.RecordFilter{
		Type:      c.Query("type"),
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
		Keyword:   c.Query("keyword"),
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			filter.CategoryID = id
		}
	}

	records, err := ctrl.repo.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "获取记录失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(records))
}

func (ctrl *RecordController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	record, err := ctrl.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "记录不存在", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(record))
}

func (ctrl *RecordController) Create(c *gin.Context) {
	var req struct {
		Type          string   `json:"type" binding:"required"`
		CategoryID    int      `json:"category_id" binding:"required"`
		SubcategoryID *int     `json:"subcategory_id"`
		Amount        float64  `json:"amount" binding:"required"`
		Date          string   `json:"date" binding:"required"`
		Note          string   `json:"note"`
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

	record := &models.Record{
		Type:          req.Type,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Amount:        req.Amount,
		Date:          req.Date,
		Note:          req.Note,
	}

	if err := ctrl.repo.Create(record); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "创建记录失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponseWithData(record, "创建成功"))
}

func (ctrl *RecordController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	var req struct {
		Type          string  `json:"type" binding:"required"`
		CategoryID    int     `json:"category_id" binding:"required"`
		SubcategoryID *int    `json:"subcategory_id"`
		Amount        float64 `json:"amount" binding:"required"`
		Date          string  `json:"date" binding:"required"`
		Note          string  `json:"note"`
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

	record, err := ctrl.repo.GetRawByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "记录不存在", "NOT_FOUND"))
		return
	}

	record.Type = req.Type
	record.CategoryID = req.CategoryID
	record.SubcategoryID = req.SubcategoryID
	record.Amount = req.Amount
	record.Date = req.Date
	record.Note = req.Note

	if err := ctrl.repo.Update(record); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "更新记录失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(record, "更新成功"))
}

func (ctrl *RecordController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseWithData(400, "无效的ID", "INVALID_ID"))
		return
	}

	if !ctrl.repo.Exists(id) {
		c.JSON(http.StatusNotFound, utils.ErrorResponseWithData(404, "记录不存在", "NOT_FOUND"))
		return
	}

	if err := ctrl.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "删除记录失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(nil, "删除成功"))
}

func (ctrl *RecordController) GetStatistics(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	stats, err := ctrl.repo.GetStatistics(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseWithData(500, "获取统计失败", "DB_ERROR"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponseWithData(stats))
}
