package jobqueue

import (
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"os"
)

func CreateWorker(src string) (w worker.Worker, err error) {
	w = worker.New()

	if err := w.Create(context.Background(), "build", "bash /build.bash"); err != nil {
		return w, err
	}

	if err := copyFileToContainer(w, "./script/build.bash", "/"); err != nil {
		return w, err
	}

	if err := copyAppCodeToWorker(w, src); err != nil {
		return w, err
	}

	return w, nil
}

func copyAppCodeToWorker(w worker.Worker, tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return w.CopyToWorker(context.Background(), r, "/app")
}

func copyFileToContainer(w worker.Worker, src string, dst string) error {
	r, err := util.ArchiveFile(src)
	if err != nil {
		return err
	}
	return w.CopyToWorker(context.Background(), r, dst)
}
