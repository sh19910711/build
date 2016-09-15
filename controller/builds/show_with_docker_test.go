// +build integration

package builds_test

import (
	"github.com/codestand/build/controller/builds"
	"github.com/codestand/build/job"
	"github.com/codestand/build/model/build"
	"github.com/codestand/build/test/testhelper/controller_helper"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestShow(t *testing.T) {
	r := gin.Default()
	builds.Mount(r)
	s := httptest.NewServer(r)
	defer s.Close()

	b := build.Build{Id: "id-foobar", Job: job.Job{Id: "myjob", Finished: true, ExitCode: 0}}
	build.Save(b)

	if res, err := controller_helper.Show(s.URL, "id-foobar"); err != nil {
		t.Fatal(err)
	} else if res.Id != b.Id {
		t.Fatal("res.Id should equal j.Id")
	} else if !res.Job.Finished {
		t.Fatal("the build job should be finished")
	} else if res.Job.ExitCode != 0 {
		t.Fatal("exit code should be zero")
	}
}
