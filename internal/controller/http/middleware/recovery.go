package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func buildPanicMessage(ctx *gin.Context, err interface{}) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.RequestURI)
	result.WriteString(" PANIC DETECTED: ")
	result.WriteString(fmt.Sprintf("%v\n%s\n", err, debug.Stack()))

	return result.String()
}

func RecoveryMiddleware(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error(buildPanicMessage(ctx, err))
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}
