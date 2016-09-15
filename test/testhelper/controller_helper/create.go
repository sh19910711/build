package controller_helper

import (
	"encoding/json"
	"errors"
	"github.com/codestand/build/test/testhelper"
	"net/http"
)

type CreateResponse struct {
	Id string
}

func Create(rootUrl, tarPath, callbackUrl string) (r CreateResponse, err error) {
	params := map[string]string{
		"callback": callbackUrl,
	}

	req, err := testhelper.UploadRequest(rootUrl+"/builds", "file", tarPath, params)
	if err != nil {
		return r, err
	}

	// send request
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return r, errors.New(res.Status)
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&r); err != nil {
		return r, err
	}
	if r.Id == "" {
		return r, errors.New("id should not be empty")
	}

	return r, nil
}
