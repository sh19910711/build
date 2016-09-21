package jobqueue

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/jobmanager"
	"github.com/codestand/build/model/job"
	"golang.org/x/net/context"
	"time"
)

var queue chan job.Job

func init() {
	queue = make(chan job.Job, 1)
}

func Queue() chan job.Job {
	return queue
}

func Close() {
	close(queue)
}

func Push(j job.Job) {
	queue <- j
	log.Debug("Push: ", j)
}

func Wait() {
	for {
		if j, ok := <-queue; ok {
			if err := spawn(j); err != nil { // TODO: async build
				log.Warn(err)
			}
		} else {
			break
		}
	}
}

func spawn(j job.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	m := jobmanager.New(ctx, j)

	if err := m.Create(j.Src); err != nil {
		return err
	}

	if err := m.Attach(); err != nil {
		log.Warn(err)
	}

	if exitCode, err := m.Run(j.Callback); err != nil {
		log.Warn(err)
		j.ExitCode = -1
	} else {
		j.ExitCode = exitCode
	}
	j.Finished = true
	job.Save(j)

	return nil
}
