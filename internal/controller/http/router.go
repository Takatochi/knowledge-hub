package http

import (
	"net/http"

	"KnowledgeHub/config"
	// Swagger documentation
	_ "KnowledgeHub/docs"
	"KnowledgeHub/internal/controller/http/middleware"
	v1 "KnowledgeHub/internal/controller/http/v1"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       KnowledgeHub API
// @description API for Knowledge Hub application
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(engine *gin.Engine, cfg *config.Config, l logger.Interface) {
	// Middleware
	engine.Use(middleware.Logger(l))
	engine.Use(middleware.Recovery(l))

	//// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// K8s health probe
	engine.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	// API v1 group
	v1Group := engine.Group("/v1")
	{
		v1.NewTranslationRoutes(v1Group, l)
	}
}
