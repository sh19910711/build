package job

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

var newJobCh chan string
var runningCh chan int

func init() {
	newJobCh = make(chan string)
	runningCh = make(chan int, 2)
}

func Wait() error {
	go watchQueue()

	for {
		log.Debug("jobqueue#wait: wait a build request")
		select {
		case buildId := <-newJobCh:
			runningCh <- 1
			log.Debug("jobqueue#wait: spawn build(" + buildId + ")")
			go spawnNewJob(buildId)

		case <-finishedCh:
			log.Debug("jobqueue#wait: finished")
			return nil
		}
	}
}

func watchQueue() {
	log.Debug("jobqueue#watch: started")
	for {
		if buildId, err := Pop(); err != nil {
			log.Warn(err)
		} else {
			log.Debug("jobqueue#watch: build(" + buildId + ") was arrived")
			newJobCh <- buildId
		}

		select {
		case <-finishedCh:
			log.Debug("jobqueue#watch: stop watching queue")
			return
		default:
		}
	}
}

func spawnNewJob(buildId string) error {
	defer func() { <-runningCh }()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	if id, err := strconv.ParseInt(buildId, 10, 32); err != nil {
		return err
	} else {
		j := New(ctx, id)
		if err := j.Create(); err != nil {
			log.Warn(err)
			return err
		}
		log.Debug("jobqueue#spawn: build(" + buildId + "): worker was created")

		if err := j.Attach(); err != nil {
			log.Warn(err)
			return err
		}
		log.Debug("jobqueue#spawn: build(" + buildId + "): worker was attached")

		if _, err := j.Run(); err != nil {
			log.Warn(err)
			return err
		}
		log.Debug("jobqueue#spawn: build(" + buildId + "): worker was finished")
	}

	return nil
}
