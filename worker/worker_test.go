package worker_test

import (
	"github.com/codestand/build/test/helper"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"testing"
)

func init() {
	helper.Init()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	w, _ := worker.New()
	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		t.Fatal(err)
	}
}
