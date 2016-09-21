package jobmanager

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func (m *JobManager) Attach() error {
	// create writer
	if err := os.MkdirAll(filepath.Dir(m.j.LogPath), 0755); err != nil {
		return err
	}

	// wait output from worker
	go func() {
		out, err := os.Create(m.j.LogPath) // TODO: use memory first?
		defer out.Close()
		if err != nil {
			log.Warn("Attach: ", err)
			return
		}

		print("attach: created: ")
		println(m.j.LogPath)
		if err := m.w.Attach(m.ctx, out); err != nil {
			log.Warn("Attach: ", err)
			return
		}
		println("attach: finished")
	}()

	return nil
}
