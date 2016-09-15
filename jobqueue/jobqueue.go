package jobqueue

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/job"
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
			if w, err := CreateWorker(j.Src); err != nil {
				log.Warn(err)
			} else {
				if exitCode, err := RunWorker(w, j.Callback); err != nil {
					log.Warn(err)
					j.ExitCode = -1
				} else {
					j.ExitCode = exitCode
				}
				j.Finished = true
				job.Save(j)
			}
		} else {
			break
		}
	}
}
