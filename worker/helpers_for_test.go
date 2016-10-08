package worker_test

import (
	"bytes"
	"github.com/codestand/build/util"
	"golang.org/x/net/context"
	"io"
	"time"
)

func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func getFakeBuildScriptTar() (io.Reader, error) {
	script := bytes.NewBufferString(`
#!/bin/sh

echo hello 1 > /dev/stdout
sleep 1
echo hello 2 > /dev/stderr
sleep 1
echo hello 3 > /dev/stdout
`)
	return util.ArchiveBuffer(script, "build.sh")
}
