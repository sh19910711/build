package jobqueue

import (
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
)

func RunWorker(w worker.Worker, callbackUrl string) (exitCode int, err error) {
	ctx := context.TODO() // with timeout?

	if err := w.Start(ctx); err != nil {
		return -1, err
	}

	if exitCode, err = w.Wait(ctx); err != nil {
		return -2, err
	}

	if err := fireCallback(w, callbackUrl); err != nil {
		return -3, err
	}

	return exitCode, nil
}

func fireCallback(w worker.Worker, callbackUrl string) error {
	if callbackUrl != "" {
		artifacts, err := w.CopyFromWorker(context.Background(), "/app/app")
		if err != nil {
			return err
		}
		if _, err := util.Upload(artifacts, callbackUrl, "file", "artifacts.tar"); err != nil {
			return err
		}
	}
	return nil
}
