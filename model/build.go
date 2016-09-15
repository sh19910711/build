package model

import (
	"github.com/codestand/build/job"
	"github.com/codestand/build/jobqueue"
	"io"
	"os"
	"path/filepath"
)

type Build struct {
	Job job.Job
}

func NewBuild() Build {
	return Build{Job: job.New()}
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

func (b *Build) SaveJob() {
	job.Save(b.Job)
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
