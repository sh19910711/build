package job

import (
	"errors"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	Id        string
	Src       string
	Callback  string
	Image     string
	Commands  []string
	Artifacts []string
	WorkerId  string
	ExitCode  int
}

func New() Job {
	return Job{Id: uuid.NewV4().String()}
}

var jobs map[string]Job

func init() {
	jobs = map[string]Job{}
}

func Save(j Job) {
	jobs[j.Id] = j
}

func Find(id string) (Job, error) {
	if j, ok := jobs[id]; ok {
		return j, nil
	} else {
		return j, errors.New("job not found " + id)
	}
}
