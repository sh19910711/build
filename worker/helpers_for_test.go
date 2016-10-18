package worker_test

import (
	"bytes"
	"github.com/codestand/build/archive"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"io"
	"strings"
	"time"
)

func contains(str, sub string) bool {
	return strings.Contains(str, sub)
}

func readerToString(r io.Reader) string {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	return b.String()
}

func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func getFakeBuildScriptTar() (io.Reader, error) {
	bs := bytes.NewBufferString(`
#!/bin/sh

echo hello 1 > /dev/stdout
sleep 1
echo hello 2 > /dev/stderr
sleep 1
echo hello 3 > /dev/stdout
`)
	return archive.TarFromBuffer(bs, "build.sh").Reader()
}

func createFakeWorker(ctx context.Context) (*worker.Worker, error) {
	w := worker.New()
	w.Image = "cs-build/test/fake-worker"
	w.Cmd = []string{"sh", "/build.sh"}

	dockerfile := bytes.NewBufferString(`
FROM alpine:3.4
CMD ["sh", "/build.sh"]
`)

	if err := w.BuildImage(ctx, dockerfile); err != nil {
		return nil, err
	}

	if err := w.Create(ctx); err != nil {
		return nil, err
	}

	if buildScriptTar, err := getFakeBuildScriptTar(); err != nil {
		return nil, err
	} else {
		if err := w.CopyToWorker(ctx, buildScriptTar, "/"); err != nil {
			return nil, err
		}
	}

	return w, nil
}
