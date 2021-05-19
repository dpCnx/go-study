package middleware

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLog() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		requestBody, _ := io.ReadAll(c.Request.Body)

		response, ok := c.Get("response")
		if !ok {
			response = ""
		}

		cost := time.Since(start).Microseconds() // 微秒
		zap.L().Info("request-->",
			zap.String("path", c.Request.URL.Path),
			zap.Any("head", c.Request.Header),
			zap.String("method", c.Request.Method),
			zap.String("body", string(requestBody)),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("response", response.(string)),
			zap.Int64("cost", cost),
		)
	}
}
