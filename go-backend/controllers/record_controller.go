package controllers

import (
	"accounting-app/data"
	"accounting-app/models"
	"accounting-app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecordController 记账记录控制器
type RecordController struct {
	recordData data.RecordData
}

// NewRecordController 创建记账记录控制器
func NewRecordController(recordData data.RecordData) *RecordController {
	return &RecordController{recordData: recordData}
}

// GetAllRecords 获取所有记账记录
func (rc *RecordController) GetAllRecords(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("categoryId"), 10, 32)

	query := &models.RecordQuery{
		Type:       c.Query("type"),
		StartDate:  c.Query("startDate"),
		EndDate:    c.Query("endDate"),
		CategoryID: uint(categoryID),
		Keyword:    c.Query("keyword"),
	}

	records, err := rc.recordData.GetAllRecords(query)
	if err != nil {
		utils.JSONDBError(c, "获取记录失败")
		return
	}

	utils.JSONSuccess(c, records)
}

// GetRecordByID 获取单个记账记录
func (rc *RecordController) GetRecordByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的记录ID")
		return
	}

	record, err := rc.recordData.GetRecordByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取记录失败")
		return
	}

	if record == nil {
		utils.JSONNotFound(c, "记录不存在")
		return
	}

	utils.JSONSuccess(c, record)
}

// CreateRecord 创建记账记录
func (rc *RecordController) CreateRecord(c *gin.Context) {
	var req struct {
		Type          string  `json:"type" binding:"required"`
		CategoryID    uint    `json:"category_id" binding:"required"`
		SubcategoryID *uint   `json:"subcategory_id"`
		Amount        float64 `json:"amount" binding:"required"`
		Date          string  `json:"date" binding:"required"`
		Note          string  `json:"note"`
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

	record := &models.Record{
		Type:          req.Type,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Amount:        req.Amount,
		Date:          req.Date,
		Note:          req.Note,
	}

	if err := rc.recordData.CreateRecord(record); err != nil {
		utils.JSONDBError(c, "创建记录失败")
		return
	}

	// 获取完整记录信息（包含关联数据）
	createdRecord, err := rc.recordData.GetRecordByID(record.ID)
	if err != nil {
		utils.JSONDBError(c, "获取新记录失败")
		return
	}

	utils.JSONCreated(c, createdRecord, "创建成功")
}

// UpdateRecord 更新记账记录
func (rc *RecordController) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的记录ID")
		return
	}

	var req struct {
		Type          string  `json:"type" binding:"required"`
		CategoryID    uint    `json:"category_id" binding:"required"`
		SubcategoryID *uint   `json:"subcategory_id"`
		Amount        float64 `json:"amount" binding:"required"`
		Date          string  `json:"date" binding:"required"`
		Note          string  `json:"note"`
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

	// 检查记录是否存在
	existingRecord, err := rc.recordData.GetRecordByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取记录失败")
		return
	}
	if existingRecord == nil {
		utils.JSONNotFound(c, "记录不存在")
		return
	}

	record := &models.Record{
		ID:            uint(id),
		Type:          req.Type,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Amount:        req.Amount,
		Date:          req.Date,
		Note:          req.Note,
	}

	if err := rc.recordData.UpdateRecord(record); err != nil {
		utils.JSONDBError(c, "更新记录失败")
		return
	}

	// 获取更新后的记录
	updatedRecord, err := rc.recordData.GetRecordByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取记录失败")
		return
	}

	utils.JSONSuccess(c, updatedRecord, "更新成功")
}

// DeleteRecord 删除记账记录
func (rc *RecordController) DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.JSONBadRequest(c, "无效的记录ID")
		return
	}

	// 检查记录是否存在
	existingRecord, err := rc.recordData.GetRecordByID(uint(id))
	if err != nil {
		utils.JSONDBError(c, "获取记录失败")
		return
	}
	if existingRecord == nil {
		utils.JSONNotFound(c, "记录不存在")
		return
	}

	if err := rc.recordData.DeleteRecord(uint(id)); err != nil {
		utils.JSONDBError(c, "删除记录失败")
		return
	}

	utils.JSONSuccess(c, nil, "删除成功")
}

// GetStatistics 获取统计信息
func (rc *RecordController) GetStatistics(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	stats, err := rc.recordData.GetStatistics(startDate, endDate)
	if err != nil {
		utils.JSONDBError(c, "获取统计失败")
		return
	}

	utils.JSONSuccess(c, stats)
}
