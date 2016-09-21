package testhelper

import (
	"bytes"
	"net/http"
)

func Get(url string, params map[string]string) (*http.Request, error) {
	return http.NewRequest("GET", url, &bytes.Buffer{})
}
