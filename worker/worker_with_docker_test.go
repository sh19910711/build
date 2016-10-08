// +build integration

package worker_test

import (
	"bytes"
	_ "github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	w := worker.New()
	w.Image = "build"
	w.Cmd = []string{"bash", "/build.bash"}
	if err := w.Create(ctx); err != nil {
		t.Fatal(err)
	}
}

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

func TestAttach(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w := worker.New()
	w.Image = "cs-build/test/attach"
	w.Cmd = []string{"sh", "/build.sh"}

	// prepare build image
	dockerfile := bytes.NewBufferString(`
FROM alpine:3.4
CMD ["sh", "/build.sh"]
`)

	if err := w.ImageBuild(ctx, dockerfile); err != nil {
		t.Fatal(err)
	}

	// create worker
	if err := w.Create(ctx); err != nil {
		t.Fatal(err)
	}

	buildScript := bytes.NewBufferString(`
#!/bin/sh

echo hello 1 > /dev/stdout
sleep 1
echo hello 2 > /dev/stderr
sleep 1
echo hello 3 > /dev/stdout
`)
	buildTar, err := util.ArchiveBuffer(buildScript, "build.sh")
	if err != nil {
		t.Fatal(err)
	}
	if err := w.CopyToWorker(ctx, buildTar, "/"); err != nil {
		t.Fatal(err)
	}

	const LOGFILE = "tmp/attach.log"

	finished := make(chan bool)
	r, err := w.Attach(ctx)
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create(LOGFILE)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		io.Copy(out, r)

		b, err := ioutil.ReadFile(LOGFILE)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "hello 1") {
			finished <- false
			return
		}
		if !strings.Contains(string(b), "hello 2") {
			finished <- false
			return
		}
		if !strings.Contains(string(b), "hello 3") {
			finished <- false
			return
		}
		finished <- true
	}()

	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if exitCode, err := w.Wait(ctx); err != nil {
		t.Fatal(err)
	} else if exitCode != 0 {
		t.Fatal("worker exited with status ", exitCode)
	}

	if ok := <-finished; !ok {
		t.Fatal("something went wrong")
	}
}
