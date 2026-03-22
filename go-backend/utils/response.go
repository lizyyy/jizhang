package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
// 与原 Node.js 项目保持 100% 格式对齐
type Response struct {
	Errno     int         `json:"Errno"`
	Errmsg    string      `json:"Errmsg"`
	ErrorCode string      `json:"ErrorCode"`
	Data      interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, message ...string) *Response {
	msg := "success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &Response{
		Errno:     0,
		Errmsg:    msg,
		ErrorCode: "OK",
		Data:      data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(errno int, errmsg string, errorCode string) *Response {
	return &Response{
		Errno:     errno,
		Errmsg:    errmsg,
		ErrorCode: errorCode,
		Data:      nil,
	}
}

// JSONSuccess 返回成功 JSON 响应（HTTP 200）
func JSONSuccess(c *gin.Context, data interface{}, message ...string) {
	c.JSON(200, SuccessResponse(data, message...))
}

// JSONCreated 返回创建成功响应（HTTP 201）
func JSONCreated(c *gin.Context, data interface{}, message ...string) {
	c.JSON(201, SuccessResponse(data, message...))
}

// JSONError 返回错误 JSON 响应
func JSONError(c *gin.Context, errno int, errmsg string, errorCode string) {
	statusCode := errno
	if statusCode < 100 || statusCode > 599 {
		statusCode = 500
	}
	c.JSON(statusCode, ErrorResponse(errno, errmsg, errorCode))
}

// JSONBadRequest 返回 400 错误
func JSONBadRequest(c *gin.Context, message string) {
	c.JSON(400, ErrorResponse(400, message, "PARAM_ERROR"))
}

// JSONNotFound 返回 404 错误
func JSONNotFound(c *gin.Context, message string) {
	c.JSON(404, ErrorResponse(404, message, "NOT_FOUND"))
}

// JSONInternalError 返回 500 错误
func JSONInternalError(c *gin.Context, message string) {
	c.JSON(500, ErrorResponse(500, message, "INTERNAL_ERROR"))
}

// JSONDBError 返回数据库错误
func JSONDBError(c *gin.Context, message string) {
	c.JSON(500, ErrorResponse(500, message, "DB_ERROR"))
}
