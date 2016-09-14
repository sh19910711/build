package jobqueue

import (
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	Id        string
	Src       string
	Callback  string
	Image     string
	Commands  []string
	Artifacts []string
}

func NewJob() Job {
	return Job{Id: uuid.NewV4().String()}
}
