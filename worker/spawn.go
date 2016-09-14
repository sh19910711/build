package worker

import (
	"github.com/codestand/build/job"
	"github.com/codestand/build/util"
	"golang.org/x/net/context"
	"os"
)

func Spawn(j job.Job) error {
	ctx := context.TODO() // with timeout?

	w, err := New()
	if err != nil {
		return err
	}

	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		return err
	}

	j.WorkerId = w.Id
	job.Save(j)

	if err := w.copyFileToContainer(ctx, "./script/build.bash", "/"); err != nil {
		return err
	}

	if err := w.copyTarBall(ctx, j.Src); err != nil {
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

func (w *Worker) copyTarBall(ctx context.Context, tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return w.CopyToWorker(ctx, r, "/app")
}

func (w *Worker) copyFileToContainer(ctx context.Context, src string, dst string) error {
	r, err := archive(src)
	if err != nil {
		return err
	}
	return w.CopyToWorker(ctx, r, dst)
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
