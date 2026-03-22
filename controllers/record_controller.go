package controllers

import (
	"accounting-app/data"
	"accounting-app/models"
	"accounting-app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecordController 记账记录控制器
type RecordController struct {
	recordData data.RecordData
}

// NewRecordController 创建记账记录控制器实例
func NewRecordController() *RecordController {
	return &RecordController{
		recordData: data.NewRecordData(),
	}
}

// GetAllRecords 获取所有记账记录
func (r *RecordController) GetAllRecords(ctx *gin.Context) {
	var filters data.RecordFilters

	filters.Type = ctx.Query("type")
	filters.StartDate = ctx.Query("startDate")
	filters.EndDate = ctx.Query("endDate")
	filters.Keyword = ctx.Query("keyword")

	if categoryIDStr := ctx.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			filters.CategoryID = categoryID
		}
	}

	records, err := r.recordData.GetAllRecords(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "获取记录失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(records, ""))
}

// GetRecordByID 根据ID获取记账记录
func (r *RecordController) GetRecordByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	record, err := r.recordData.GetRecordByID(id)
	if err != nil {
		if err.Error() == "记录不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "记录不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "获取记录失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(record, ""))
}

// CreateRecord 创建记账记录
func (r *RecordController) CreateRecord(ctx *gin.Context) {
	var record models.Record
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验必要参数
	if record.Type == "" || record.CategoryID == 0 || record.Amount == 0 || record.Date == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	// 校验类型
	validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
	if !validTypes[record.Type] {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的类型", "INVALID_TYPE"))
		return
	}

	created, err := r.recordData.CreateRecord(&record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "创建记录失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(created, "创建成功"))
}

// UpdateRecord 更新记账记录
func (r *RecordController) UpdateRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	var record models.Record
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error(), "PARAM_ERROR"))
		return
	}

	// 校验必要参数
	if record.Type == "" || record.CategoryID == 0 || record.Amount == 0 || record.Date == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "缺少必要参数", "PARAM_ERROR"))
		return
	}

	// 校验类型
	validTypes := map[string]bool{"expense": true, "income": true, "transfer": true, "loan": true}
	if !validTypes[record.Type] {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的类型", "INVALID_TYPE"))
		return
	}

	record.ID = id

	updated, err := r.recordData.UpdateRecord(&record)
	if err != nil {
		if err.Error() == "记录不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "记录不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "更新记录失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(updated, "更新成功"))
}

// DeleteRecord 删除记账记录
func (r *RecordController) DeleteRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的ID", "PARAM_ERROR"))
		return
	}

	err = r.recordData.DeleteRecord(id)
	if err != nil {
		if err.Error() == "记录不存在" {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(404, "记录不存在", "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "删除记录失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "删除成功"))
}

// GetStatistics 获取统计信息
func (r *RecordController) GetStatistics(ctx *gin.Context) {
	var filters data.StatisticsFilters
	filters.StartDate = ctx.Query("startDate")
	filters.EndDate = ctx.Query("endDate")

	stats, err := r.recordData.GetStatistics(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "获取统计失败", "DB_ERROR"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(stats, ""))
}
