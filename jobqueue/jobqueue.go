package jobqueue

import (
	log "github.com/Sirupsen/logrus"
)

var queue chan Job

func init() {
	queue = make(chan Job, 1)
}

func Queue() chan Job {
	return queue
}

func Close() {
	close(queue)
}

func Push(j Job) {
	queue <- j
	log.Debug("Push: ", j)
}

func Wait() {
	for {
		if j, ok := <-queue; ok {
			spawnJob(j) // TODO: parallelize
		} else {
			break
		}
	}
}
