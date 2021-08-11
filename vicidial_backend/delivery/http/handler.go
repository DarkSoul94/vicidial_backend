package http

import (
	"io/ioutil"
	"net/http"

	"github.com/DarkSoul94/vicidial_backend/helper"
	helperUC "github.com/DarkSoul94/vicidial_backend/helper/usecase"
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

	data := make(map[string]interface{})
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

	data := make(map[string]interface{})
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
		return
	}

	data = h.prepareData(data, "inn", "phone")

	data = map[string]interface{}{
		"flag_get": "get_info_vici",
		"inn":      data["inn"],
		"phone":    data["phone"],
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

	params := h.getParamsFromUrl(c)
	params = h.prepareData(params, "phone", "inn", "send_sms")
	data := map[string]interface{}{
		"phone":    params["phone"],
		"inn":      params["inn"],
		"send_sms": params["send_sms"],
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

	data := make(map[string]interface{})
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
	}

	data = h.prepareData(data, "phone", "inn", "send_sms")
	data = map[string]interface{}{
		"phone":    data["phone"],
		"inn":      data["inn"],
		"send_sms": data["send_sms"],
	}
	response, _ := h.makeRequestTo1c("ivr", data)

	c.JSON(http.StatusOK, response)
}

//getParamsFromUrl получает все параметры которые были переданы и формирует
func (h *Handler) getParamsFromUrl(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	keys := make([]string, 0, len(c.Request.URL.Query()))
	for k := range c.Request.URL.Query() {
		keys = append(keys, k)
	}
	for _, val := range keys {
		params[val] = c.Request.URL.Query().Get(val)
	}
	return params
}

//prepareData функция заполняет пустой строкой значение в мапе если значение по ключу не было найдено
func (h *Handler) prepareData(data map[string]interface{}, keys ...string) map[string]interface{} {
	for _, val := range keys {
		if _, ok := data[val]; !ok {
			data[val] = ""
		}
	}
	return data
}
