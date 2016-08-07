package worker_test

import (
	"github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"testing"
)

func init() {
	testhelper.Init()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	w, _ := worker.New()
	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		t.Fatal(err)
	}
}
