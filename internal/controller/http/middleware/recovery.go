package middleware

import (
	"KnowledgeHub/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strings"
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

// Видаляємо невикористану функцію logPanic
// func logPanic(l logger.LoggerInterface) func(ctx *gin.Context, err interface{}) {
//     return func(ctx *gin.Context, err interface{}) {
//         l.Error(buildPanicMessage(ctx, err))
//     }
// }

func Recovery(l logger.LoggerInterface) gin.HandlerFunc {
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
