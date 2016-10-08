// +build integration

package worker_test

import (
	"bytes"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"testing"
	"time"
)

const DOCKERFILE_WITHOUT_CMD string = `
FROM alpine:3.4
RUN echo hello
`

const DOCKERFILE_FAILED string = `
FROM alpine:3.4
RUN false
`

func TestImageBuild(t *testing.T) {
	t.Run("without cmd", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		w := worker.New()
		w.Image = "cs-build/test/hello"
		buf := bytes.NewBufferString(DOCKERFILE_WITHOUT_CMD)
		if err := w.ImageBuild(ctx, buf); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		w := worker.New()
		w.Image = "cs-build/test/fail"
		buf := bytes.NewBufferString(DOCKERFILE_FAILED)
		if err := w.ImageBuild(ctx, buf); err == nil {
			t.Fatal("build should be failed")
		}
	})
}
