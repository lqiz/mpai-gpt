package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func WithTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), "trace_id", strings.ReplaceAll(traceId, "-", ""))
		c.Request = c.Request.WithContext(ctx)
		c.Header("trace-id", traceId)
	}
}
