package http

import (
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

// HelloWorld ...
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

	//TODO form request in UC
	//request := h.uc.MakeRequestTo1C(data)

	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetLKInfo(c *gin.Context) {
	var err error

	if err = h.validateAuthKey(c); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	if err = c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": ErrDataIsNotJson.Error()})
	}
	/*
		data = map[string]interface{}{
			"flag_get": "get_info_vici",
			"inn":      data["inn"],
			"phone":    data["phone"],
		}

		gtToken := viper.GetString("app.auth_getaway_token")
		gtURL := viper.GetString("app.getaway_url")

		temp, _ := json.Marshal(data)
		body := bytes.NewBuffer(temp)
	*/
	//c.JSON(http.StatusOK, resp)
}
