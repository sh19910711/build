package job

import (
	"errors"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	Id        string   `json:"id"`
	Src       string   `json:"src"`
	Callback  string   `json:"callback"`
	Image     string   `json:"image"`
	Commands  []string `json:"commands"`
	Artifacts []string `json:"artifacts"`
	WorkerId  string   `json:"worker_id"`
	ExitCode  int      `json:"exit_code"`
	Finished  bool     `json:"finished"`
}

func New() Job {
	return Job{Id: uuid.NewV4().String()}
}

var jobs map[string]Job // TODO: let's use database

func init() {
	jobs = map[string]Job{}
}

func Save(j Job) error {
	jobs[j.Id] = j
	return nil
}

func Find(id string) (Job, error) {
	if j, ok := jobs[id]; ok {
		return j, nil
	} else {
		return j, errors.New("job not found " + id)
	}
}
