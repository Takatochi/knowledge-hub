package middleware

import (
	"strconv"
	"strings"

	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func buildRequestMessage(ctx *gin.Context, status int, bodySize int) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.RequestURI)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(status))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(bodySize))

	return result.String()
}

func LoggerMiddleware(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Перед виконанням основного хендлера
		ctx.Next()

		// Після того як обробка завершена
		status := ctx.Writer.Status()
		bodySize := ctx.Writer.Size()

		l.Info(buildRequestMessage(ctx, status, bodySize))
	}
}
