package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HttpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		logger := log.Info()
		status := c.Writer.Status()
		statusCode := http.StatusText(c.Writer.Status())
		if status != http.StatusTemporaryRedirect {
			logger = log.Error()
		}

		logger.Str("protocol", "http").
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Int("status_code", status).
			Str("status_text", statusCode).
			Dur("duration", duration).
			Msg("received a HTTP request")
	}
}
