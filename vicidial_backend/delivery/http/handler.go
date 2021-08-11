package http

import (
	"net/http"

	"github.com/DarkSoul94/vicidial_backend/helper"
	helperUC "github.com/DarkSoul94/vicidial_backend/helper/usecase"
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
)

// Handler ...
type Handler struct {
	uc         vicidial_backend.Usecase
	httpClient helper.Helper
}

// NewHandler ...
func NewHandler(uc vicidial_backend.Usecase) *Handler {
	return &Handler{
		uc:         uc,
		httpClient: helperUC.NewHelper(),
	}
}

// HelloWorld ...
func (h *Handler) HelloWorld(c *gin.Context) {
	h.uc.HelloWorld(c.Request.Context())
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
