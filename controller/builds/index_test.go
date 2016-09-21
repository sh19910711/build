package builds_test

import (
	"github.com/codestand/build/controller/builds"
	"github.com/codestand/build/model/build"
	"github.com/codestand/build/test/testhelper/controller_helper"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	r := gin.Default()
	builds.Mount(r)
	s := httptest.NewServer(r)
	defer s.Close()

	b := build.Build{Id: "hello"}
	build.Save(b)

	if res, err := controller_helper.Index(s.URL); err != nil {
		t.Fatal(err)
	} else {
		if res.Builds[0].Id != b.Id {
			t.Fatal("/builds returns all build")
		}
	}
}
