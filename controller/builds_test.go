package controller_test

import (
	"encoding/json"
	"github.com/codestand/build/controller"
	"github.com/codestand/build/controller/test/helper"
	"github.com/codestand/build/jobqueue"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	testhelper.Init()
	jobqueue.Init()

	// start server
	s := testhelper.Serve("/builds", controller.Create)
	defer s.Close()

	// prepare request
	params := map[string]string{
		"callback": s.URL + "/callback",
	}
	req, err := testhelper.UploadRequest(s.URL+"/builds", "file", "../example/app.tar", params)
	if err != nil {
		t.Fatal(err)
	}

	// send request
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// check response
	dec := json.NewDecoder(res.Body)
	type Output struct {
		Id string
	}

	// check output
	var o Output
	if err := dec.Decode(&o); err != nil {
		t.Fatal(err)
	}
	if o.Id == "" {
		t.Fatal(o)
	}
	res.Body.Close()

	// check status
	if res.StatusCode != 200 {
		t.Fatal(res.StatusCode)
	}
}
