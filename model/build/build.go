package build

import (
	"errors"
	"github.com/codestand/build/jobqueue"
	"github.com/codestand/build/model/job"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path/filepath"
)

type Build struct {
	Id    string  `json:"id"`
	Job   job.Job `json:"job"` // TODO: has many
	JobId string  `json:"-"`
}

var builds map[string]Build

func init() {
	builds = map[string]Build{}
}

func New() Build {
	j := job.New()
	return Build{Id: uuid.NewV4().String(), Job: j, JobId: j.Id}
}

func Find(id string) (b Build, err error) {
	if b, ok := builds[id]; ok {
		if b.JobId != "" {
			if j, err := job.Find(b.JobId); err != nil {
				return b, errors.New("the build job is not found")
			} else {
				b.Job = j
			}
		}
		return b, nil
	}
	return b, errors.New("the build is not found")
}

func Save(b Build) error {
	if err := job.Save(b.Job); err != nil {
		return err
	}
	builds[b.Id] = b
	return nil
}

func (b *Build) SetCallbackURL(url string) {
	b.Job.Callback = url
}

func (b *Build) SetWorker() { // TODO: use config
	b.Job.Image = "build"
	b.Job.Commands = []string{"make"}
	b.Job.Artifacts = []string{"/app/app"}
}

func (b *Build) PushJobQueue() {
	jobqueue.Push(b.Job)
}

func (b *Build) SaveSourceCode(r io.Reader, prefix string) error {
	if err := os.MkdirAll(prefix, 0700); err != nil {
		return err
	}

	src := filepath.Join(prefix, b.Job.Id)
	if err := writeToFile(r, src); err != nil {
		return err
	} else {
		b.Job.Src = src
	}

	return nil
}

func writeToFile(r io.Reader, path string) error {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}
