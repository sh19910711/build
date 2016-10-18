package job

import (
	"github.com/codestand/build/model"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
)

type Job struct {
	ctx context.Context
	w   *worker.Worker
	B   *model.Build
}

func New(ctx context.Context, b *model.Build) *Job {
	return &Job{
		ctx: ctx,
		B:   b,
	}
}
