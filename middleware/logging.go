package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志中间件，记录请求输入、输出和错误
func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 获取请求信息
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		clientIP := ctx.ClientIP()

		// 读取请求体
		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录请求
		log.Printf("[REQUEST] %s %s | IP: %s | Body: %s", method, path, clientIP, string(requestBody))

		// 使用自定义响应写入器捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		// 处理请求
		ctx.Next()

		// 计算耗时
		latency := time.Since(startTime)

		// 获取状态码
		statusCode := ctx.Writer.Status()

		// 获取响应体
		responseBody := blw.body.String()

		// 检查是否有错误
		if len(ctx.Errors) > 0 {
			log.Printf("[ERROR] %s %s | Status: %d | Errors: %v", method, path, statusCode, ctx.Errors)
		}

		// 记录响应
		log.Printf("[RESPONSE] %s %s | Status: %d | Latency: %v | Body: %s",
			method, path, statusCode, latency, responseBody)
	}
}

// bodyLogWriter 用于捕获响应体的写入器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w bodyLogWriter) WriteHeaderNow() {
	w.ResponseWriter.WriteHeaderNow()
}

// JSONLogMiddleware 结构化JSON日志中间件
func JSONLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var requestBody map[string]interface{}
		if ctx.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &requestBody)
		}

		// 处理请求
		ctx.Next()

		// 计算耗时
		latency := time.Since(startTime)

		// 构建日志条目
		logEntry := map[string]interface{}{
			"timestamp":  startTime.Format(time.RFC3339),
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"client_ip":  ctx.ClientIP(),
			"status":     ctx.Writer.Status(),
			"latency":    latency.Milliseconds(),
			"user_agent": ctx.Request.UserAgent(),
			"request":    requestBody,
		}

		// 如果有错误，添加到日志中
		if len(ctx.Errors) > 0 {
			logEntry["errors"] = ctx.Errors.Errors()
		}

		// 输出JSON日志
		logJSON, _ := json.Marshal(logEntry)
		log.Println(string(logJSON))
	}
}
