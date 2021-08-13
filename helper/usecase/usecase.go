package usecase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

type Helper struct {
	client *http.Client
}

func NewHelper() *Helper {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.ResponseHeaderTimeout = 30 * time.Second
	return &Helper{
		client: &http.Client{
			Transport: tr,
			Timeout:   60 * time.Second,
		},
	}
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

	responce, err = h.client.Do(request)
	if err != nil {
		return &http.Response{}, err
	}

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

	responce, err = h.client.Do(request)
	if err != nil {
		return &http.Response{}, err
	}

	return responce, nil
}
