package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/job"
	"github.com/codestand/build/util"
	"golang.org/x/net/context"
	"os"
)

func Spawn(j job.Job) error {
	log.Debug("worker: spawn job: ", j)
	ctx := context.TODO() // with timeout?

	w, err := New()
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

func (w *Worker) sendTarBall(ctx context.Context, tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return w.Copy(ctx, r, "/app")
}

func (w *Worker) fireCallback(ctx context.Context, j job.Job) error {
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
