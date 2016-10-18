package job

import (
	"github.com/codestand/build/worker"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
)

func (j *Job) Create() error {
	j.w = worker.New()

	if hasDockerFile, r, err := getDockerfileFromSource(j); err != nil {
		return err
	} else if hasDockerFile {
		if err := buildWithDockerfileIfExists(j); err != nil {
			return err
		}
	}

	if err := j.w.Create(j.ctx); err != nil {
		return err
	}

	if err := copyFileToContainer(j, "./script/build.bash", "/"); err != nil {
		return err
	}

	if err := copyAppCodeToWorker(j, src); err != nil {
		return err
	}

	return nil
}

func getDockerfileFromTar(j *Job) (ok bool, nilReader io.Reader, err error) {
	var filename string = "Dockerfile"
	tarball, err := os.Open(tarPath)
	if err != nil {
		return false, nilReader, err
	}
	defer tarball.Close()

	if ok, err := util.CheckFileInTar(tarball, filename); err != nil {
		return false, nilReader, err
	} else if ok {
		if r, err := util.ReadFileFromTar(tarball, "Dockerfile"); err != nil {
			return false, nilReader, err
		} else {
			return true, r, nil
		}
	} else {
		return false, nilReader, nil
	}
}

func buildWithDockerfileIfExists(j *Job, r io.Reader) error {
	j.w.Image = uuid.NewV4().String()
	if err := j.w.ImageBuild(j.ctx, r); err != nil {
		return err
	} else {
		return nil
	}
}

func copyAppCodeToWorker(j *Job, tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return j.w.CopyToWorker(j.ctx, r, "/app")
}

func copyFileToContainer(j *Job, src string, dst string) error {
	r, err := util.ArchiveFile(src)
	if err != nil {
		return err
	}
	return j.w.CopyToWorker(j.ctx, r, dst)
}
