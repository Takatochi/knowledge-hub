package v1

import (
	_ "KnowledgeHub/docs"
	"KnowledgeHub/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// V1 -.
type V1 struct {
	l logger.LoggerInterface
	v *validator.Validate
}
type entity struct {
	Message string `json:"message" example:"success"`
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} entity
// @Failure     500 {object} response.Error
// @Router      /translation/history [get]
func (r *V1) history(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"data": entity{"Server is running"},
	})
}
