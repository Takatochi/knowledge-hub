package v1

import (
	"net/http"

	// Swagger documentation
	_ "KnowledgeHub/docs"
	"KnowledgeHub/internal/models"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	l logger.Interface
	v *validator.Validate
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Entity
// @Failure     500 {object} response.Error
// @Router      /translation/history [get]
func (r *V1) history(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": models.Entity{Message: "Server is running"},
	})
}
