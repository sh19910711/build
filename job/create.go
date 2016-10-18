package job

import (
	"bytes"
	"fmt"
	"github.com/codestand/build/archive"
	"github.com/codestand/build/worker"
)

func (j *Job) Create() error {
	j.w = worker.New()
	j.w.Image = "build"
	j.w.Cmd = []string{"bash", "/build.bash"}

	if err := j.w.Create(j.ctx); err != nil {
		return err
	}

	if err := copyBuildScriptToContainer(j); err != nil {
		return err
	}

	if err := copyAppCodeToWorker(j); err != nil {
		return err
	}

	return nil
}

func copyAppCodeToWorker(j *Job) error {
	if app, err := j.B.AppTar(); err != nil {
		return err
	} else {
		return j.w.CopyToWorker(j.ctx, app, "/app")
	}
}

func copyBuildScriptToContainer(j *Job) error {
	script := ""
	script += fmt.Sprintln("#!/bin/bash")
	script += fmt.Sprintln("make")

	r, err := archive.TarFromBuffer(bytes.NewBufferString(script), "build.bash").Reader()
	if err != nil {
		return err
	}
	return j.w.CopyToWorker(j.ctx, r, "/")
}
