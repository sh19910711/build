package jobqueue

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/job"
	"github.com/codestand/build/worker"
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
			if err := worker.Spawn(j); err != nil { // TODO: parallelize
				log.Fatal(err) // TODO: improve error handling
			}
		} else {
			break
		}
	}
}
