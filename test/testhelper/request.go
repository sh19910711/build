package testhelper

import (
	"bytes"
	"net/http"
)

func Get(url string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	return http.NewRequest("GET", url, body)
}
