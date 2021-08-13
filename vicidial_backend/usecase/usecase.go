package usecase

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/DarkSoul94/vicidial_backend/helper"
	"github.com/DarkSoul94/vicidial_backend/models"
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

func (u *Usecase) AddLeads(leads []models.Lead) {
	go u.addLeads(leads)
}

func (u *Usecase) addLeads(leads []models.Lead) {
	for _, lead := range leads {
		u.addLead(lead)
	}
}

func (u *Usecase) addLead(lead models.Lead) {
	var body []byte
	const SuccessLeadAdd string = "SUCCESS: add_lead LEAD HAS BEEN ADDED"

	max_tries := viper.GetInt("app.vicidial.max_tries")
	resource := "/vicidial/non_agent_api.php"
	_, tz := time.Now().Local().Zone()

	data := map[string]interface{}{
		"phone_number":    lead.Get("phone_number", ""),
		"list_id":         lead.Get("list_id", ""),
		"security_phrase": lead.Get("security_phrase", ""),
		"address1":        lead.Get("address1", ""),
		"address2":        lead.Get("address2", ""),
		"address3":        lead.Get("address3", ""),
		"province":        lead.Get("province", ""),
		"last_name":       lead.Get("last_name", ""),
		"postal_code":     lead.Get("postal_code", ""),
		"city":            lead.Get("city", ""),
		"email":           lead.Get("email", ""),
		"first_name":      lead.Get("first_name", ""),
		"phone_code":      lead.Get("phone_code", ""),
		"source":          "test",
		"user":            viper.GetString("app.vicidial.login"),
		"pass":            viper.GetString("app.vicidial.pass"),
		"gmt_offset_now":  fmt.Sprint(time.Duration(tz * int(time.Second)).Hours()), //
		"function":        "add_lead",
	}

	if _, ok := lead["include_lists"]; ok {
		data["action"] = "add_unique_lead"
		data["include_lists"] = lead.Get("include_lists", "")
		data["exclude_statuses"] = lead.Get("exclude_statuses", "")
		resource = "/non_agent_api_ext/index.php"
	}

	if _, ok := lead["callback"]; ok {
		data["callback"] = "Y"
		data["callback_status"] = "CALLBK"
		data["campaign_id"] = "ccMain"
		data["callback_comments"] = lead.Get("callback_comments", "")
		data["callback_datetime"] = lead.Get("callback_datetime", lead.Get("security_phrase", ""))
	}
	url := viper.GetString("app.vicidial.url") + resource

	success := false
	for i := 0; i < max_tries; i++ {
		res, err := u.httpClient.Get(url, data)
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed GET-request to %s", resource), "add lead", data["phone_number"].(string), err)
			time.Sleep(viper.GetDuration("app.vicidial.delay") * time.Millisecond)
			continue
		}

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed read response body from %s", url), "add lead", data["action"].(string), err)
		}
		res.Body.Close()
		if res.StatusCode == 200 || strings.Contains(string(body), SuccessLeadAdd) || res.Status == "200 OK" {
			success = true
			break
		}

		time.Sleep(viper.GetDuration("app.vicidial.delay") * time.Millisecond)
	}
	if !success {
		logger.LogError(fmt.Sprintf("Failed add lead to %s", resource), "add lead", data["phone_number"].(string), errors.New(string(body)))
	}
}

func (u *Usecase) UpdateLead(lead models.Lead) {
	resource := "/vicidial/non_agent_api.php"
	url := viper.GetString("app.vicidial.url") + resource

	delete(lead, "type")

	lead["function"] = "update_lead"
	lead["user"] = viper.GetString("app.vicidial.login")
	lead["pass"] = viper.GetString("app.vicidial.pass")
	lead["source"] = "test"

	res, err := u.httpClient.Get(url, lead)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed GET-request to %s", resource), "update lead", resource, err)
		return
	}

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		logger.LogError("Failed update lead", "update lead", string(body), nil)
	}

	res.Body.Close()

}
