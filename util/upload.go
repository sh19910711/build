package util

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func Upload(r io.Reader, url, key, filename string) (*http.Response, error) {
	req, err := uploadRequest(r, url, key, filename)
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	return c.Do(req)
}

func uploadRequest(r io.Reader, url string, key string, filename string) (req *http.Request, err error) {
	// prepare writer
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	// add form file
	part, err := w.CreateFormFile(key, filename)
	if err != nil {
		return nil, err
	}
	io.Copy(part, r)
	if err := w.Close(); err != nil {
		return nil, err
	}

	// create http.Request
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}
