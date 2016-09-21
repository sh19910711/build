package controller_helper

import (
	"encoding/json"
	"errors"
	"github.com/codestand/build/test/testhelper"
	"net/http"
)

type IndexResponse struct {
	Builds []BuildResponse `json:"builds"`
}

func Index(url string) (b IndexResponse, err error) {
	req, err := testhelper.Get(url+"/builds", map[string]string{})
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
