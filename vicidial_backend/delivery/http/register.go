package http

import (
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints ...
func RegisterHTTPEndpoints(router *gin.Engine, uc vicidial_backend.Usecase) {
	h := NewHandler(uc)

	apiEndpoints := router.Group("/api")
	{
		apiEndpoints.POST("/", h.HelloWorld)
	}
}
