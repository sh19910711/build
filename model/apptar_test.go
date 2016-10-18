package model_test

import (
	"bytes"
	"github.com/codestand/build/archive"
	"github.com/codestand/build/model"
	"strings"
	"testing"
)

func TestAppTar(t *testing.T) {
	model.Open()
	defer model.Close()

	b := &model.Build{Id: 10000}
	if model.Find(b).RecordNotFound() {
		t.Fatal("id=10000 should be found")
	}

	if apptar, err := b.AppTar(); err != nil {
		t.Fatal(err)
	} else if r, err := archive.GetFileReaderFromTar(apptar, "main.c"); err != nil {
		t.Fatal(err)
	} else {
		b := &bytes.Buffer{}
		b.ReadFrom(r)
		if !strings.Contains(b.String(), "stdio.h") {
			t.Fatal("main.c should contain stdio.h")
		}
	}
}
