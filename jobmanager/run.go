package jobmanager

import (
	"github.com/codestand/build/util"
)

func (m *JobManager) Run(callbackUrl string) (exitCode int, err error) {
	if err := m.w.Start(m.ctx); err != nil {
		return -1, err
	}

	if exitCode, err = m.w.Wait(m.ctx); err != nil {
		return -2, err
	}

	if err := m.fireCallback(callbackUrl); err != nil {
		return -3, err
	}

	return exitCode, nil
}

func (m *JobManager) fireCallback(callbackUrl string) error {
	if callbackUrl != "" {
		artifacts, err := m.w.CopyFromWorker(m.ctx, "/app/app")
		if err != nil {
			return err
		}
		// TODO: use context
		if _, err := util.Upload(artifacts, callbackUrl, "file", "artifacts.tar"); err != nil {
			return err
		}
	}
	return nil
}
