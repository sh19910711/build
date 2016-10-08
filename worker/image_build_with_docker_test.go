// +build integration

package worker_test

import (
	"bytes"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestImageBuild(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	w := worker.New()

	t.Run("hello", func(t *testing.T) {
		w.Image = "cs-build/test/hello"
		buf := bytes.NewBufferString(`
FROM alpine:3.4
RUN echo hello
`)
		if err := w.ImageBuild(ctx, buf); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("fail on run command", func(t *testing.T) {
		w.Image = "cs-build/test/fail"
		buf := bytes.NewBufferString(`
FROM alpine:3.4
RUN false
`)
		if err := w.ImageBuild(ctx, buf); err == nil {
			t.Fatal("build should be failed")
		}
	})
}
