package v1

import (
	// Swagger documentation
	_ "KnowledgeHub/docs"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group *gin.RouterGroup, l logger.Interface) {
	r := &V1{l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
	}
}
