package model_test

import (
	"github.com/codestand/build/model"
	"strings"
	"testing"
)

func TestAllModels(t *testing.T) {
	model.Open()
	defer model.Close()

	builds := model.All()
	for _, b := range builds {
		if b.Id == 10000 {
			return
		}
	}
	t.Fatal("id=10000 should be found")
}

func TestModelFound(t *testing.T) {
	model.Open()
	defer model.Close()

	b := model.Build{Id: 10000}
	model.Find(&b)
	if b.Id != 10000 {
		t.Fatal("id=10000 should be found")
	}
}

func TestModelNotFound(t *testing.T) {
	model.Open()
	defer model.Close()

	b := model.Build{Id: 999999}
	if !model.Find(&b).RecordNotFound() {
		t.Fatal("id=999999 should not be found")
	}
}

func TestWriteLog(t *testing.T) {
	model.Open()
	defer model.Close()

	b := model.Build{Id: 10000}
	model.Find(&b)
	b.WriteLog("hello")
	b.WriteLog("world")

	model.Find(&b)
	if !strings.Contains(b.Log, "hello\nworld") {
		t.Fatal("the log should contain hello and world")
	}
}
