package jobqueue

import (
	uuid "github.com/satori/go.uuid"
)

type job struct {
	Id        string
	Src       string
	Callback  string
	Image     string
	Commands  []string
	Artifacts []string
}

func NewJob() job {
	u := uuid.NewV4()
	return job{Id: u.String()}
}
