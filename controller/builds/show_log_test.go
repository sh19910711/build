package builds_test

import (
	"fmt"
	"github.com/codestand/build/controller/builds"
	"github.com/codestand/build/model/build"
	"github.com/codestand/build/model/job"
	"github.com/codestand/build/test/testhelper"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShowLog(t *testing.T) {
	// prepare server
	r := gin.Default()
	builds.Mount(r)
	s := httptest.NewServer(r)
	defer s.Close()

	// prepare temp dir
	d, err := ioutil.TempDir("", "test-show-log")
	if err != nil {
		t.Fatal(err)
	}

	// save build
	b := build.Build{
		Id: "build-id",
		Job: job.Job{
			Id:      "job-id",
			LogPath: filepath.Join(d, "log.txt"),
		},
	}
	build.Save(b)

	// prepare fake log file
	if w, err := os.Create(b.Job.LogPath); err != nil {
		t.Fatal(err)
	} else {
		fmt.Fprintln(w, "hello world")
		w.Close()
	}

	// request
	req, err := testhelper.Get(s.URL+"/builds/build-id/log.txt", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// check response
	if b, err := ioutil.ReadAll(res.Body); err != nil {
		t.Fatal(err)
	} else if !strings.Contains(string(b), "hello world") {
		t.Fatal("response should contain hello world")
	}
}
