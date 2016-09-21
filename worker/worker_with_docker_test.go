// +build integration

package worker_test

import (
	"bytes"
	_ "github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	w := worker.New()
	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		t.Fatal(err)
	}
}

func TestImageBuild(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	w := worker.New()

	t.Run("hello world", func(t *testing.T) {
		buf := bytes.NewBufferString(`
FROM alpine:3.4
RUN echo hello
`)
		if err := w.ImageBuild(ctx, "hello", buf); err != nil {
			t.Fatal(err)
		}
	})
}
