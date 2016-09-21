package jobmanager

import (
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	uuid "github.com/satori/go.uuid"
	"os"
)

func (m *JobManager) Create(src string) (err error) {
	m.w = worker.New()

	if err := m.buildWithDockerfileIfExists(src); err != nil {
		return err
	}

	if err := m.w.Create(m.ctx, m.w.ImageName, "bash /build.bash"); err != nil {
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

func (m *JobManager) buildWithDockerfileIfExists(tarPath string) error {
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
			if err := m.w.ImageBuild(m.ctx, imageName, r); err != nil {
				return err
			} else {
				m.w.ImageName = imageName
			}
		}
	}
	return nil
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
