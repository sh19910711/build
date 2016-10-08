// +build integration

package worker_test

import (
	_ "github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/worker"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	w := worker.New()
	w.Image = "build"
	w.Cmd = []string{"bash", "/build.bash"}

	if err := w.Create(ctx); err != nil {
		t.Fatal(err)
	}
}
