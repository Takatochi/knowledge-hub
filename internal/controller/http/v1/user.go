package v1

import (
	"net/http"
	"strconv"

	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
	logger      logger.Interface
}

func NewUserHandler(userService *services.UserService, logger logger.Interface) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get user details by ID
// @ID           get-user-by-id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
