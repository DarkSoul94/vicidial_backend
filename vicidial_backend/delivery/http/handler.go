package http

import (
	"io/ioutil"
	"net/http"

	"github.com/DarkSoul94/vicidial_backend/helper"
	helperUC "github.com/DarkSoul94/vicidial_backend/helper/usecase"
	"github.com/DarkSoul94/vicidial_backend/models"
	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	allowedActions []string = []string{
		"send_sms", "get_payment_requsits", "get_main_info",
		"get_loan_info", "get_ticket_info", "get_detail_loan",
		"get_detail_ticket", "find_by_fio", "get_loans_by_phone",
		"get_balance_on_date", "get_phones_from_order",
	}
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

//validateAuthKey берет из Header`a авторизационный ключ и сравнивает с тем что указан в конфиг файле
func (h *Handler) validateAuthKey(c *gin.Context) error {
	authKey := c.GetHeader("X-Auth-Key")
	if authKey != viper.GetString("app.auth_key") {
		return ErrNotAuthenticated
	}
	return nil
}

// VicidialActions ...
func (h *Handler) VicidialActions(c *gin.Context) {
	var (
		err error
	)

	if err = h.validateAuthKey(c); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	action := c.Param("action")

	actionAllowed := func(action string) bool {
		for _, val := range allowedActions {
			if action == val {
				return true
			}
		}
		return false
	}

	if !actionAllowed(action) {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrMethodNotAlowed.Error()})
		return
	}

	data := make(models.Lead)
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
		return
	}
	data["action"] = action

	response, _ := h.makeRequestTo1c("vicidial", data)

	c.JSON(http.StatusOK, response)
}

// GetLKInfo ...
func (h *Handler) GetLKInfo(c *gin.Context) {
	var err error

	if err = h.validateAuthKey(c); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	data := make(models.Lead)
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
		return
	}

	data = map[string]interface{}{
		"flag_get": "get_info_vici",
		"inn":      data.Get("inn", ""),
		"phone":    data.Get("phone", ""),
	}
	gtUrl := viper.GetString("app.getaway_url")
	gtToken := viper.GetString("app.auth_getaway_token")

	objResponse, err := h.httpClient.Post(
		gtUrl,
		data,
		map[string]string{"token": gtToken})

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	body, _ := ioutil.ReadAll(objResponse.Body)
	c.JSON(http.StatusOK, body)
}

func (h *Handler) IvrGet(c *gin.Context) {
	var err error

	if err = h.validateAuthKey(c); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	lead := h.getLeadFromUrl(c)
	data := map[string]interface{}{
		"phone":    lead.Get("phone", ""),
		"inn":      lead.Get("inn", ""),
		"send_sms": lead.Get("send_sms", false),
	}
	response, _ := h.makeRequestTo1c("ivr", data)

	c.JSON(http.StatusOK, response)

}

func (h *Handler) IvrPost(c *gin.Context) {
	var err error

	if err = h.validateAuthKey(c); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	data := make(models.Lead)
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
		return
	}

	data = map[string]interface{}{
		"phone":    data.Get("phone", ""),
		"inn":      data.Get("inn", ""),
		"send_sms": data.Get("send_sms", false),
	}
	response, _ := h.makeRequestTo1c("ivr", data)

	c.JSON(http.StatusOK, response)
}

//getParamsFromUrl получает все параметры которые были переданы и формирует
func (h *Handler) getLeadFromUrl(c *gin.Context) models.Lead {
	params := make(models.Lead)
	keys := make([]string, 0, len(c.Request.URL.Query()))
	for k := range c.Request.URL.Query() {
		keys = append(keys, k)
	}
	for _, val := range keys {
		params[val] = c.Request.URL.Query().Get(val)
	}
	return params
}

func (h *Handler) AddLead(c *gin.Context) {
	var (
		data []models.Lead = make([]models.Lead, 0)
		err  error
	)

	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
		return
	}

	h.uc.AddLeads(data)

	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateLead(c *gin.Context) {
	var (
		data map[string]interface{} = make(map[string]interface{})
		err  error
	)

	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
		return
	}

	h.uc.UpdateLead(data)

	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) NonAgentApi(c *gin.Context) {
	var (
		data     map[string]interface{}
		resource string
		result   map[string]string
		err      error
	)

	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
		return
	}

	resource = "/vicidial/non_agent_api.php"
	result, err = h.makeRequestToVicidial(resource, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
