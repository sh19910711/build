package model

import (
	"github.com/codestand/build/job"
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

func (b *Build) Save(r io.Reader, prefix string) error {
	if err := os.MkdirAll(prefix, 0700); err != nil {
		return err
	}

	src := filepath.Join(prefix, b.Job.Id)
	if err := save(r, src); err != nil {
		return err
	} else {
		b.Job.Src = src
	}

	return nil
}

func save(r io.Reader, path string) error {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}
