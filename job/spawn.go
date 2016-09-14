package job

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"os"
)

func (j *Job) Spawn() error {
	log.Debug("jobqueue: spawn: ", j)
	ctx := context.TODO() // with timeout?

	w, err := newBuildWorker()
	if err != nil {
		return err
	}

	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		return err
	}

	if err := w.CopyFile(ctx, "./script/build.bash", "/"); err != nil {
		return err
	}

	if err := w.sendTarBall(ctx, j.Src); err != nil {
		return err
	}

	if err := w.Start(ctx); err != nil {
		return err
	}

	if exitCode, err := w.Wait(ctx); err != nil {
		return err
	} else {
		j.ExitCode = exitCode
	}

	if err := w.fireCallback(ctx, j); err != nil {
		return err
	}

	return nil
}

type buildWorker struct {
	*worker.Worker
}

func newBuildWorker() (buildWorker, error) {
	if w, err := worker.New(); err != nil {
		return buildWorker{}, err
	} else {
		return buildWorker{&w}, nil
	}
}

func (w *buildWorker) sendTarBall(ctx context.Context, tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return w.Copy(ctx, r, "/app")
}

func (w *buildWorker) fireCallback(ctx context.Context, j *Job) error {
	if j.Callback != "" {
		artifacts, err := w.CopyFromWorker(ctx, "/app/app")
		if err != nil {
			return err
		}
		if _, err := util.Upload(artifacts, j.Callback, "file", "artifacts.tar"); err != nil {
			return err
		}
	}
	return nil
}
