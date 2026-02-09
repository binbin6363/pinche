package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"pinche/internal/logger"
)

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// read request body
		var requestBody string
		if c.Request.Body != nil && c.Request.ContentLength > 0 && c.Request.ContentLength < 10240 {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = string(bodyBytes)
			// mask sensitive fields in request body
			requestBody = maskSensitiveData(requestBody)
		}

		// wrap response writer
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = rw

		// process request
		c.Next()

		// calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// get user_id if authenticated
		userID, _ := c.Get("user_id")

		// response body (truncate if too large)
		responseBody := rw.body.String()
		if len(responseBody) > 2048 {
			responseBody = responseBody[:2048] + "...[truncated]"
		}

		// log fields
		fields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"latency_ms", latency.Milliseconds(),
			"client_ip", clientIP,
		}

		if query != "" {
			fields = append(fields, "query", query)
		}
		if userID != nil {
			fields = append(fields, "user_id", userID)
		}
		if requestBody != "" && requestBody != "{}" {
			fields = append(fields, "request_body", requestBody)
		}
		if statusCode >= 400 {
			fields = append(fields, "response_body", responseBody)
		}
		if userAgent != "" {
			fields = append(fields, "user_agent", userAgent)
		}

		// log based on status code
		if statusCode >= 500 {
			logger.Error("HTTP Request", fields...)
		} else if statusCode >= 400 {
			logger.Warn("HTTP Request", fields...)
		} else {
			logger.Info("HTTP Request", fields...)
		}
	}
}

// maskSensitiveData masks sensitive fields in JSON string
func maskSensitiveData(data string) string {
	// simple masking for common sensitive fields
	// in production, consider using proper JSON parsing
	sensitiveFields := []string{"password", "token", "secret"}
	result := data
	for _, field := range sensitiveFields {
		// mask values like "password":"xxx" or "password": "xxx"
		result = maskJSONField(result, field)
	}
	return result
}

// maskJSONField masks a specific field value in JSON string
func maskJSONField(data, field string) string {
	// simple string-based masking
	patterns := []string{
		`"` + field + `":"`,
		`"` + field + `": "`,
	}
	for _, pattern := range patterns {
		idx := bytes.Index([]byte(data), []byte(pattern))
		if idx == -1 {
			continue
		}
		start := idx + len(pattern)
		end := start
		for end < len(data) && data[end] != '"' {
			end++
		}
		if end < len(data) {
			// replace the value with **** and include closing quote
			data = data[:start] + "****" + data[end:]
		}
	}
	return data
}
