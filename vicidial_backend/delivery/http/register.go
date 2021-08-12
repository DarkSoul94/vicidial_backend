package http

import (
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints ...
func RegisterHTTPEndpoints(router *gin.RouterGroup, uc vicidial_backend.Usecase) {
	h := NewHandler(uc)

	vicidialEndpoints := router.Group("/vicidial")
	{
		vicidialEndpoints.POST("/:action", h.VicidialActions)
	}

	ivrEndpoint := router.Group("/ivr")
	{
		// /api/ivr
		ivrEndpoint.GET("", h.IvrGet)
		ivrEndpoint.POST("", h.IvrPost)

	}

	// /api/add_lead
	router.Group("/add_lead").POST("", h.AddLead)
	// /api/update_lead
	router.Group("/update_lead").POST("", h.UpdateLead)
	// /api/non_agent_api
	router.Group("/non_agent_api").POST("", h.NonAgentApi)

}
