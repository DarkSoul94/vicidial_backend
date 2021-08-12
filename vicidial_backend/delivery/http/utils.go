package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/DarkSoul94/vicidial_backend/models"
	"github.com/spf13/viper"
)

func (h *Handler) makeRequestTo1c(resource string, data map[string]interface{}) (map[string]interface{}, error) {
	var (
		url    string
		header map[string]string      = make(map[string]string)
		result map[string]interface{} = make(map[string]interface{})
	)

	url = viper.GetString("app.1C.url") + resource

	header["Authorization"] = fmt.Sprintf("Basic %s", viper.GetString("app.1C.auth_token"))

	res, err := h.httpClient.Post(url, data, header)
	if err != nil {
		return map[string]interface{}{"error": "connection error"}, nil
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	json.Unmarshal(body, &result)

	return result, nil
}

func (h *Handler) makeRequestToVicidial(resource string, data map[string]interface{}) (map[string]interface{}, error) {
	var (
		url      string
		req_data map[string]interface{} = make(map[string]interface{})
		result   map[string]interface{} = make(map[string]interface{})
	)

	url = viper.GetString("app.vicidial.url") + resource

	req_data["user"] = viper.GetString("app.vicidial.login")
	req_data["pass"] = viper.GetString("app.vicidial.pass")
	req_data["source"] = "test"

	res, err := h.httpClient.Get(url, req_data)
	if err != nil {
		return map[string]interface{}{}, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	json.Unmarshal(body, &result)

	return result, nil
}

func (h *Handler) makeRequestToGateway(data models.Lead) (map[string]interface{}, error) {
	response := make(map[string]interface{})
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
		return nil, err
	}

	body, _ := ioutil.ReadAll(objResponse.Body)
	objResponse.Body.Close()
	json.Unmarshal(body, &response)
	return response, nil
}
