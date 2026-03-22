package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = rw

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		var requestBodyFormatted interface{}
		if len(requestBody) > 0 {
			var jsonBody map[string]interface{}
			if err := json.Unmarshal(requestBody, &jsonBody); err == nil {
				requestBodyFormatted = jsonBody
			} else {
				requestBodyFormatted = string(requestBody)
			}
		}

		var responseBodyFormatted interface{}
		responseBody := rw.body.String()
		if len(responseBody) > 0 {
			var jsonBody map[string]interface{}
			if err := json.Unmarshal([]byte(responseBody), &jsonBody); err == nil {
				responseBodyFormatted = jsonBody
			} else {
				responseBodyFormatted = responseBody
			}
		}

		logEntry := map[string]interface{}{
			"timestamp":   start.Format(time.RFC3339),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"query":       c.Request.URL.RawQuery,
			"status_code": statusCode,
			"duration":    duration.String(),
			"client_ip":   c.ClientIP(),
			"request":     requestBodyFormatted,
			"response":    responseBodyFormatted,
		}

		if len(c.Errors) > 0 {
			logEntry["errors"] = c.Errors.String()
		}

		logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
		println(string(logJSON))
	}
}
