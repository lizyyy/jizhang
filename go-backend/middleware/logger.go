package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 自定义响应写入器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggerMiddleware 日志中间件
// 记录每一次请求的输入、输出和错误
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// 构建日志信息
		logEntry := fmt.Sprintf("\n================== REQUEST LOG ==================\n")
		logEntry += fmt.Sprintf("Time:       %s\n", time.Now().Format("2006-01-02 15:04:05"))
		logEntry += fmt.Sprintf("Method:     %s\n", method)
		logEntry += fmt.Sprintf("Path:       %s\n", path)
		logEntry += fmt.Sprintf("Client IP:  %s\n", clientIP)
		logEntry += fmt.Sprintf("Status:     %d\n", statusCode)
		logEntry += fmt.Sprintf("Duration:   %v\n", duration)

		// 记录查询参数
		if query := c.Request.URL.RawQuery; query != "" {
			logEntry += fmt.Sprintf("Query:      %s\n", query)
		}

		// 记录请求体（限制长度）
		if len(requestBody) > 0 {
			bodyStr := string(requestBody)
			if len(bodyStr) > 1000 {
				bodyStr = bodyStr[:1000] + "..."
			}
			logEntry += fmt.Sprintf("Request:    %s\n", bodyStr)
		}

		// 记录响应体（限制长度）
		responseBody := writer.body.String()
		if len(responseBody) > 0 {
			if len(responseBody) > 1000 {
				responseBody = responseBody[:1000] + "..."
			}
			logEntry += fmt.Sprintf("Response:   %s\n", responseBody)
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			logEntry += fmt.Sprintf("Errors:     %v\n", c.Errors)
		}

		logEntry += "================================================\n"

		// 输出日志
		fmt.Print(logEntry)
	}
}
