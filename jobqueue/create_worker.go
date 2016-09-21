package jobqueue

import (
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"os"
)

func CreateWorker(src string) (w worker.Worker, err error) {
	w = worker.New()

	if err := buildWithDockerfileIfExists(w, src); err != nil {
		return w, err
	}

	if err := w.Create(context.Background(), w.ImageName, "bash /build.bash"); err != nil {
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

func buildWithDockerfileIfExists(w worker.Worker, tarPath string) error {
	tarball, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	if ok, err := util.CheckFileInTar(tarball, "Dockerfile"); err != nil {
		return err
	} else if ok {
		if r, err := util.ReadFileFromTar(tarball, "Dockerfile"); err != nil {
			return err
		} else {
			imageName := uuid.NewV4().String()
			if err := w.ImageBuild(context.Background(), imageName, r); err != nil {
				return err
			} else {
				w.ImageName = imageName
			}
		}
	}
	return nil
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
