// +build docker

package builds_test

import (
	"encoding/json"
	"errors"
	"github.com/codestand/build/controller/builds"
	"github.com/codestand/build/jobqueue"
	"github.com/codestand/build/test/testhelper"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	//
	// set up fake server
	//
	r := gin.Default()
	builds.Mount(r)

	// test artifacts send to callback URL
	r.POST("/callback", func(c *gin.Context) {
		r, _, err := c.Request.FormFile("file")
		if err != nil {
			t.Error(err)
		}
		res, err := testhelper.ShouldIncludeFileInTar(r, "app")
		if err != nil {
			t.Error(err)
		}
		if !res {
			t.Error("artifact should be found")
		}
	})

	s := httptest.NewServer(r)
	defer s.Close()

	// prepare jobqueue
	go jobqueue.Wait()
	defer jobqueue.Close()

	// prepare request
	params := map[string]string{
		"callback": s.URL + "/callback",
	}
	req, err := testhelper.UploadRequest(s.URL+"/builds", "file", "./example/app.tar", params)
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
	buildId, err := checkCreateResponse(res)
	if err != nil {
		t.Fatal(err)
	}

	// wait for finishing build
	exitCode := make(chan int, 1)
	go func() {
		for {
			if res, err := getBuild(s.URL, buildId); err != nil {
				t.Fatal(err)
			} else {
				if res.Finished {
					exitCode <- res.ExitCode
					break
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// timeout after three seconds
	select {
	case c := <-exitCode:
		if c != 0 {
			t.Fatal(res)
		}
	case <-time.After(3 * time.Second):
		t.Fatal(errors.New("build timed out"))
	}
}

type BuildResponse struct {
	Id       string
	ExitCode int
	Finished bool
}

func getBuild(url, id string) (b BuildResponse, err error) {
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

type CreatedBuildResponse struct {
	Id string
}

func checkCreateResponse(res *http.Response) (string, error) {
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New(res.Status)
	}

	dec := json.NewDecoder(res.Body)
	var b CreatedBuildResponse
	if err := dec.Decode(&b); err != nil {
		return "", err
	}
	if b.Id == "" {
		return "", errors.New("id should not be empty")
	}

	return b.Id, nil
}
