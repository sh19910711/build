package jobqueue

import (
	log "github.com/Sirupsen/logrus"
)

var queue chan job

func Queue() chan job {
	return queue
}

func Init() {
	queue = make(chan job, 1)
}

func Close() {
	close(queue)
}

func Push(newjob job) {
	go func() {
		queue <- newjob
		log.Debug("jobqueue: pushed: ", newjob)
	}()
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
