package jobmanager

import (
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
)

func (m *JobManager) Create(src string) error {
	m.w = worker.New()

	if ok, r, err := m.getDockerFile(src); err != nil {
		return err
	} else if ok {
		if err := m.buildWithDockerfileIfExists(r); err != nil {
			return err
		}
	}

	if err := m.w.Create(m.ctx); err != nil {
		return err
	}

	if err := m.copyFileToContainer("./script/build.bash", "/"); err != nil {
		return err
	}

	if err := m.copyAppCodeToWorker(src); err != nil {
		return err
	}

	return nil
}

func (m *JobManager) getDockerFile(tarPath string) (ok bool, nilReader io.Reader, err error) {
	var filename string = "Dockerfile"
	tarball, err := os.Open(tarPath)
	defer tarball.Close()
	if err != nil {
		return false, nilReader, err
	}
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

func (m *JobManager) buildWithDockerfileIfExists(r io.Reader) error {
	m.w.Image = uuid.NewV4().String()
	if err := m.w.ImageBuild(m.ctx, r); err != nil {
		return err
	} else {
		return nil
	}
}

func (m *JobManager) copyAppCodeToWorker(tarPath string) error {
	r, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	return m.w.CopyToWorker(m.ctx, r, "/app")
}

func (m *JobManager) copyFileToContainer(src string, dst string) error {
	r, err := util.ArchiveFile(src)
	if err != nil {
		return err
	}
	return m.w.CopyToWorker(m.ctx, r, dst)
}
