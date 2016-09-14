package job

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
	ExitCode  int
}

func New() Job {
	return Job{Id: uuid.NewV4().String()}
}
