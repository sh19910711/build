// +build integration

package builds_test

import (
	"github.com/codestand/build/controller/builds"
	"github.com/codestand/build/jobqueue"
	"github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/test/testhelper/controller_helper"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
	"time"
)

// The Create should create new build and callback after finished build
func TestCreate(t *testing.T) {
	// set up fake web server
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

	// run web server
	s := httptest.NewServer(r)
	defer s.Close()

	// prepare jobqueue
	go jobqueue.Wait()
	defer jobqueue.Close()

	// send request
	build, err := controller_helper.Create(s.URL, "./example/app.tar", s.URL+"/callback")
	if err != nil {
		t.Fatal(err)
	}

	// wait for finishing build
	exitCode := make(chan int, 1)
	go func() {
		for {
			if res, err := controller_helper.Show(s.URL, build.Id); err != nil {
				t.Fatal(err)
			} else {
				if res.Job.Finished {
					exitCode <- res.Job.ExitCode
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
			t.Fatal(c)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the build should be finished in a few second")
	}
}
