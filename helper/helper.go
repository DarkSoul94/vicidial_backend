package helper

import "net/http"

type Helper interface {
	Get(url string, data map[string]interface{}) (*http.Response, error)
	Post(url string, data map[string]interface{}, headers map[string]string) (*http.Response, error)
}
