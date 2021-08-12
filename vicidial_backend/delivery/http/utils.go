package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	json.Unmarshal(body, &result)

	return result, nil
}

func (h *Handler) makeRequestToVicidial(resource string, data map[string]interface{}) (map[string]string, error) {
	var (
		url      string
		req_data map[string]interface{} = make(map[string]interface{})
		result   map[string]string      = map[string]string{}
	)

	url = viper.GetString("app.vicidial.url") + resource

	req_data["user"] = viper.GetString("app.vicidial.login")
	req_data["pass"] = viper.GetString("app.vicidial.pass")
	req_data["source"] = "test"

	res, err := h.httpClient.Get(url, req_data)
	if err != nil {
		return map[string]string{}, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	json.Unmarshal(body, &result)

	return result, nil
}
