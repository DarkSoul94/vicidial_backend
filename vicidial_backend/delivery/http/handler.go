package http

import (
	"net/http"

	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
)

// Handler ...
type Handler struct {
	uc vicidial_backend.Usecase
}

// NewHandler ...
func NewHandler(uc vicidial_backend.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// HelloWorld ...
func (h *Handler) HelloWorld(c *gin.Context) {
	h.uc.HelloWorld(c.Request.Context())
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
