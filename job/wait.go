package job

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

func Wait() error {
	go watchQueue()

	for {
		select {
		case buildId := <-newJobCh:
			go spawnNewJob(buildId)
		case <-finishedCh:
			return nil
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func watchQueue() {
	for {
		if buildId, err := Pop(); err != nil {
			log.Warn(err)
		} else {
			newJobCh <- buildId
		}
		select {
		case <-finishedCh:
			return
		}
	}
}

func spawnNewJob(buildId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	if id, err := strconv.ParseInt(buildId, 10, 32); err != nil {
		return err
	} else {
		j := New(ctx, id)
		if err := j.Create(); err != nil {
			return err
		}
		if err := j.Attach(); err != nil {
			return err
		}
		if _, err := j.Run(); err != nil {
			return err
		}
	}

	return nil
}
