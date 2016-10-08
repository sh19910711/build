// +build integration

package worker_test

import (
	"bytes"
	"github.com/codestand/build/worker"
	"testing"
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
		ctx, cancel := contextWithTimeout()
		defer cancel()

		w := worker.New()
		w.Image = "cs-build/test/hello"

		r := bytes.NewBufferString(DOCKERFILE_WITHOUT_CMD)
		if err := w.ImageBuild(ctx, r); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		ctx, cancel := contextWithTimeout()
		defer cancel()

		w := worker.New()
		w.Image = "cs-build/test/fail"

		r := bytes.NewBufferString(DOCKERFILE_FAILED)
		if err := w.ImageBuild(ctx, r); err == nil {
			t.Fatal("build should be failed")
		}
	})
}
