package http

import (
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints ...
func RegisterHTTPEndpoints(router *gin.Engine, uc vicidial_backend.Usecase) {
	h := NewHandler(uc)

	vicidialEndpoints := router.Group("/vicidial")
	{
		vicidialEndpoints.POST("/get_lk_info", h.GetLKInfo)
		vicidialEndpoints.POST("/:action", h.VicidialActions)
	}

	/*ivrEndpoint := router.Group("/ivr")
	{
		ivrEndpoint.GET("", h.Ivr)
		ivrEndpoint.POST("", h.Ivr)

	}
	*/

	router.Group("/add_lead").POST("", h.AddLead)
	router.Group("/update_lead").POST("", h.UpdateLead)
	router.Group("/non_agent_api").POST("", h.NonAgentApi)

}
