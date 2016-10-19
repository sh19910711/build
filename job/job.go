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

func New(ctx context.Context, buildId int64) *Job {
	b := &model.Build{Id: buildId}
	if model.Find(b).RecordNotFound() {
		return nil
	}
	return &Job{
		ctx: ctx,
		B:   b,
	}
}
