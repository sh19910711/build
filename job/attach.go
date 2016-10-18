package job

import (
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func (m *JobManager) Attach() error {
	// create writer
	if err := os.MkdirAll(filepath.Dir(m.j.LogPath), 0755); err != nil {
		return err
	}
	out, err := os.Create(m.j.LogPath) // TODO: use memory first?
	if err != nil {
		return err
	}

	r, err := m.w.Attach(m.ctx)
	if err != nil {
		return err
	}
	// wait output from worker
	go func() {
		io.Copy(out, r)
		if err := out.Close(); err != nil {
			log.Warn("Attach: ", err)
		}
	}()

	return nil
}
