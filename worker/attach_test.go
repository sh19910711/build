// +build integration

package worker_test

import (
	"bytes"
	_ "github.com/codestand/build/test/testhelper"
	"github.com/codestand/build/worker"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAttach(t *testing.T) {
	ctx, cancel := contextWithTimeout()
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

	if buildScriptTar, err := getFakeBuildScriptTar(); err != nil {
		t.Fatal(err)
	} else {
		if err := w.CopyToWorker(ctx, buildScriptTar, "/"); err != nil {
			t.Fatal(err)
		}
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
