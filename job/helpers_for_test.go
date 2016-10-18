package job_test

import (
	"github.com/codestand/build/model"
	"golang.org/x/net/context"
	"strings"
	"time"
)

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
	b := &model.Build{Id: 10000}
	model.Find(b)
	return b
}

func contains(str, sub string) bool {
	return strings.Contains(str, sub)
}
