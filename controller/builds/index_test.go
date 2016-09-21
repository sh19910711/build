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
		var ok bool = false
		for _, b := range res.Builds {
			if b.Id == b.Id {
				ok = true
			}
		}
		if !ok {
			t.Fatal("GET /builds should return all build")
		}
	}
}
