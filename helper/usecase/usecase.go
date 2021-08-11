package usecase

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) Get(url string, data map[string]interface{}) (*http.Response, error) {
	var (
		request  *http.Request
		responce *http.Response
		err      error
	)

	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	query := request.URL.Query()
	for key, val := range data {
		query.Set(key, val.(string))
	}
	request.URL.RawQuery = query.Encode()

	client := &http.Client{}
	responce, err = client.Do(request)
	if err != nil {
		return &http.Response{}, err
	}
	defer responce.Body.Close()
	return responce, nil
}

func (h *Helper) Post(url string, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	var (
		request  *http.Request
		responce *http.Response
		err      error
	)

	bt, _ := json.Marshal(&data)
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(bt))
	if err != nil {
		return &http.Response{}, err
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{}
	responce, err = client.Do(request)
	if err != nil {
		return &http.Response{}, err
	}
	defer responce.Body.Close()
	return responce, nil

}