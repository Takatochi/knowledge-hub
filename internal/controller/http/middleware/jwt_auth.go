package middleware

import (
	"net/http"

	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtService *services.JWTService, logger logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logger.Info("Missing Authorization header from %s", ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		token, err := jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			logger.Info("Invalid Authorization header format from %s: %v", ctx.ClientIP(), err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			logger.Info("Token validation failed from %s: %v", ctx.ClientIP(), err)

			var errorMessage string
			switch err {
			case services.ErrExpiredToken:
				errorMessage = "Token has expired"
			case services.ErrInvalidToken:
				errorMessage = "Invalid token"
			case services.ErrInvalidClaims:
				errorMessage = "Invalid token claims"
			default:
				errorMessage = "Token validation failed"
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errorMessage,
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("jwt_claims", claims)

		logger.Info("User %s (ID: %d) authenticated successfully from %s",
			claims.Username, claims.UserID, ctx.ClientIP())

		ctx.Next()
	}
}

// OptionalJWTAuthMiddleware створює middleware для опціональної аутентифікації
// Не блокує запит, якщо токен відсутній, але валідує його, якщо присутній
func OptionalJWTAuthMiddleware(jwtService *services.JWTService, logger logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		token, err := jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			logger.Info("Invalid Authorization header format from %s: %v", ctx.ClientIP(), err)
			ctx.Next()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			logger.Info("Token validation failed from %s: %v", ctx.ClientIP(), err)
			ctx.Next()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("jwt_claims", claims)

		logger.Info("User %s (ID: %d) optionally authenticated from %s",
			claims.Username, claims.UserID, ctx.ClientIP())

		ctx.Next()
	}
}

func GetUserIDFromContext(ctx *gin.Context) (uint, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

func GetUsernameFromContext(ctx *gin.Context) (string, bool) {
	username, exists := ctx.Get("username")
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

func GetEmailFromContext(ctx *gin.Context) (string, bool) {
	email, exists := ctx.Get("email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}

func GetJWTClaimsFromContext(ctx *gin.Context) (*services.JWTClaims, bool) {
	claims, exists := ctx.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*services.JWTClaims)
	return jwtClaims, ok
}
