package job_test

import (
	"github.com/codestand/build/model"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func sleep(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}

const FAKE_BUILD_ID = 10000

func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

func setup() {
	model.Open()
}

func teardown() {
	model.Close()
}

func getFakeBuild() *model.Build {
	b := &model.Build{Id: FAKE_BUILD_ID}
	if model.Find(b).RecordNotFound() {
		return nil
	} else {
		return b
	}
}

func contains(str, sub string) bool {
	return strings.Contains(str, sub)
}
