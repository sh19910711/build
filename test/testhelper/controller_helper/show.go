package controller_helper

import (
	"encoding/json"
	"errors"
	"github.com/codestand/build/test/testhelper"
	"net/http"
)

type ShowResponse struct {
	Id  string
	Job ShowJobResponse
}

type ShowJobResponse struct {
	ExitCode int `json:exitcode`
	Finished bool
}

func Show(url, id string) (b ShowResponse, err error) {
	req, err := testhelper.Get(url+"/builds/"+id, map[string]string{})
	if err != nil {
		return b, err
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return b, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return b, errors.New(res.Status)
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&b); err != nil {
		return b, err
	}
	return b, nil
}
