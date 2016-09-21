package jobmanager

import (
	"github.com/codestand/build/model/job"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
)

type JobManager struct {
	ctx context.Context
	w   worker.Worker
	j   job.Job
}

func New(ctx context.Context, j job.Job) JobManager {
	return JobManager{
		ctx: ctx,
		j:   j,
	}
}
