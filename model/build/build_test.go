package build_test

import (
	"github.com/codestand/build/model/build"
	_ "github.com/codestand/build/test/testhelper"
	"testing"
)

func TestSaveAndFind(t *testing.T) {
	b := build.Build{Id: "id-foobar"}
	build.Save(b)
	if found, err := build.Find("id-foobar"); err != nil {
		t.Fatal(err)
	} else if found.Id != b.Id {
		t.Fatal("build id is wrong")
	}
}

func TestSaveAndFindWithJob(t *testing.T) {
	b := build.New()
	b.Id = "id-foobar"
	build.Save(b)
	if found, err := build.Find("id-foobar"); err != nil {
		t.Fatal(err)
	} else if found.Id != b.Id {
		t.Fatal("build id is wrong")
	}
}
