package v1

import (
	"net/http"

	"KnowledgeHub/internal/controller/http/middleware"
	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler обробляє запити аутентифікації
type AuthHandler struct {
	jwtService  *services.JWTService
	userService *services.UserService
	logger      logger.Interface
	validator   *validator.Validate
}

// NewAuthHandler створює новий екземпляр AuthHandler
func NewAuthHandler(jwtService *services.JWTService, userService *services.UserService, logger logger.Interface) *AuthHandler {
	return &AuthHandler{
		jwtService:  jwtService,
		userService: userService,
		logger:      logger,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

// LoginRequest представляє запит на логін
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest представляє запит на реєстрацію
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// AuthResponse представляє відповідь після успішної аутентифікації
type AuthResponse struct {
	AccessToken  string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string   `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt    int64    `json:"expires_at" example:"1640995200"`
	User         UserInfo `json:"user"`
}

// UserInfo представляє інформацію про користувача
type UserInfo struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"johndoe@example.com"`
}

// RefreshRequest представляє запит на оновлення токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse представляє структуру помилки
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT tokens
// @ID           login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} AuthResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// TODO: Тут буде логіка перевірки користувача в базі даних
	// Поки що заглушку для демонстрації
	if req.Username != "admin" || req.Password != "password" {
		h.logger.Info("Failed login attempt for username: %s from %s", req.Username, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	tokenPair, err := h.jwtService.GenerateTokenPair(1, req.Username, "admin@example.com")
	if err != nil {
		h.logger.Error("Failed to generate tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authentication tokens",
		})
		return
	}

	h.logger.Info("User %s logged in successfully from %s", req.Username, c.ClientIP())

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: UserInfo{
			ID:       1,
			Username: req.Username,
			Email:    "admin@example.com",
		},
	})
}

// Register godoc
// @Summary      User registration
// @Description  Register a new user and return JWT tokens
// @ID           register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration data"
// @Success      201 {object} AuthResponse
// @Failure      400 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid registration request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// TODO: Тут буде логіка створення користувача в базі даних
	// Поки що заглушку для демонстрації
	if req.Username == "admin" {
		h.logger.Info("Registration attempt with existing username: %s from %s", req.Username, c.ClientIP())
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	userID := uint(2) // Заглушка для ID нового користувача
	tokenPair, err := h.jwtService.GenerateTokenPair(userID, req.Username, req.Email)
	if err != nil {
		h.logger.Error("Failed to generate tokens for new user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authentication tokens",
		})
		return
	}

	h.logger.Info("User %s registered successfully from %s", req.Username, c.ClientIP())

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: UserInfo{
			ID:       userID,
			Username: req.Username,
			Email:    req.Email,
		},
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Refresh access token using refresh token
// @ID           refresh-token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh token"
// @Success      200 {object} AuthResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Валідуємо refresh токен
	_, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Info("Invalid refresh token from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// TODO: Тут буде логіка отримання користувача з бази даних за ID з claims
	// Поки що використовуємо заглушку
	userID := uint(1)
	username := "admin"
	email := "admin@example.com"

	// Генеруємо нові токени
	tokenPair, err := h.jwtService.RefreshTokens(req.RefreshToken, userID, username, email)
	if err != nil {
		h.logger.Error("Failed to refresh tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh tokens",
		})
		return
	}

	h.logger.Info("Tokens refreshed successfully for user %s from %s", username, c.ClientIP())

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: UserInfo{
			ID:       userID,
			Username: username,
			Email:    email,
		},
	})
}

// MessageResponse представляє відповідь з повідомленням
type MessageResponse struct {
	Message string `json:"message" example:"Successfully logged out"`
}

// Logout godoc
// @Summary      User logout
// @Description  Logout user (invalidate tokens on client side)
// @ID           logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} MessageResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Отримуємо інформацію про користувача з контексту
	username, exists := middleware.GetUsernameFromContext(c)
	if !exists {
		username = "unknown"
	}

	h.logger.Info("User %s logged out from %s", username, c.ClientIP())

	// В JWT архітектурі logout зазвичай обробляється на клієнті
	// шляхом видалення токенів з локального сховища
	// Для серверного logout можна використовувати blacklist токенів
	c.JSON(http.StatusOK, MessageResponse{
		Message: "Successfully logged out",
	})
}

// Me godoc
// @Summary      Get current user info
// @Description  Get information about currently authenticated user
// @ID           me
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} UserInfo
// @Failure      401 {object} ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	username, _ := middleware.GetUsernameFromContext(c)
	email, _ := middleware.GetEmailFromContext(c)

	c.JSON(http.StatusOK, UserInfo{
		ID:       userID,
		Username: username,
		Email:    email,
	})
}
