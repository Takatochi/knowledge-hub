package v1

import (
	// Swagger documentation
	_ "KnowledgeHub/docs"
	"KnowledgeHub/internal/controller/http/middleware"
	"KnowledgeHub/internal/services"
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

func NewAuthRoutes(apiV1Group *gin.RouterGroup, jwtService *services.JWTService, userService *services.UserService, l logger.Interface) {

	authHandler := NewAuthHandler(jwtService, userService, l)
	// Публічні роути (без аутентифікації)
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// Захищені роути (з обов'язковою аутентифікацією)
	protectedAuthGroup := apiV1Group.Group("/auth")
	protectedAuthGroup.Use(middleware.JWTAuthMiddleware(jwtService, l))
	{
		protectedAuthGroup.POST("/logout", authHandler.Logout)
		protectedAuthGroup.GET("/me", authHandler.Me)
	}
}
