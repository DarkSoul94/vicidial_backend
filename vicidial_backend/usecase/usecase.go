package usecase

import (
	"fmt"
	"time"

	"github.com/DarkSoul94/vicidial_backend/helper"
	"github.com/DarkSoul94/vicidial_backend/pkg/logger"
	"github.com/spf13/viper"
)

// Usecase ...
type Usecase struct {
	httpClient helper.Helper
}

// NewUsecase ...
func NewUsecase(httpUC helper.Helper) *Usecase {
	return &Usecase{
		httpClient: httpUC,
	}
}

func (u *Usecase) AddLeads(leads []map[string]interface{}) {
	go u.addLeads(leads)
}

func (u *Usecase) addLeads(leads []map[string]interface{}) {
	for _, lead := range leads {
		u.addLead(lead)
	}
}

func (u *Usecase) addLead(lead map[string]interface{}) {
	max_tries := viper.GetInt("app.vicidial.max_tries")
	resource := "/vicidial/non_agent_api.php"
	url := viper.GetString("app.vicidial.url") + resource
	_, tz := time.Now().Local().Zone()

	data := map[string]interface{}{
		"phone_number":    getLeadParam(lead, "phone_number"),
		"list_id":         getLeadParam(lead, "list_id"),
		"security_phrase": getLeadParam(lead, "security_phrase"),
		"address1":        getLeadParam(lead, "address1"),
		"address2":        getLeadParam(lead, "address2"),
		"address3":        getLeadParam(lead, "address3"),
		"province":        getLeadParam(lead, "province"),
		"last_name":       getLeadParam(lead, "last_name"),
		"postal_code":     getLeadParam(lead, "postal_code"),
		"city":            getLeadParam(lead, "city"),
		"email":           getLeadParam(lead, "email"),
		"first_name":      getLeadParam(lead, "first_name"),
		"phone_code":      getLeadParam(lead, "phone_code"),
		"source":          getLeadParam(lead, "source"),
		"user":            viper.GetString("app.vicidial.login"),
		"pass":            viper.GetString("app.vicidial.pass"),
		"gmt_offset_now":  fmt.Sprint(time.Duration(tz * int(time.Second)).Hours()), //
		"function":        "add_lead",
	}

	if _, ok := lead["include_lists"]; ok {
		data["action"] = "add_unique_lead"
		data["include_lists"] = getLeadParam(lead, "include_lists")
		data["exclude_statuses"] = getLeadParam(lead, "exclude_statuses")
		resource = "/non_agent_api_ext/index.php"
	}

	if _, ok := lead["callback"]; ok {
		data["callback"] = "Y"
		data["callback_status"] = "CALLBK"
		data["campaign_id"] = "ccMain"
		data["callback_comments"] = getLeadParam(lead, "callback_comments")

		if val, ok := lead["callback_datetime"]; ok {
			data["callback_datetime"] = val.(string)
		} else {
			if val, ok := lead["security_phrase"]; ok {
				data["callback_datetime"] = val.(string)
			} else {
				data["callback_datetime"] = ""
			}
		}
	}

	success := false
	for i := 0; i < max_tries; i++ {
		res, err := u.httpClient.Get(url, data)
		if err != nil {
			logger.LogError("Get method failed", "add lead", data["phone_number"].(string), err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if res.StatusCode == 200 {
			success = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !success {
		logger.LogError("Failed add lead", "add lead", data["phone_number"].(string), nil)
	}
}

func getLeadParam(lead map[string]interface{}, param string) string {
	if val, ok := lead[param]; ok {
		return val.(string)
	} else {
		return ""
	}
}

func (u *Usecase) UpdateLead(lead map[string]interface{}) {
	resource := "/vicidial/non_agent_api.php"
	url := viper.GetString("app.vicidial.url") + resource

	delete(lead, "type")

	lead["function"] = "update_lead"
	lead["user"] = viper.GetString("app.vicidial.login")
	lead["pass"] = viper.GetString("app.vicidial.pass")
	lead["source"] = "test"

	u.httpClient.Get(url, lead)
}
