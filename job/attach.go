package job

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
)

func (j *Job) Attach() error {
	r, err := j.w.Attach(j.ctx)
	if err != nil {
		return err
	}

	// wait output from worker
	go j.handleWorkerOutput(r)

	return nil
}

func (j *Job) handleWorkerOutput(r io.Reader) {
	j.B.ResetLog()
	br := bufio.NewReaderSize(r, 2048)
	for {
		if line, _, err := br.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			log.Warn("ERROR: Attach: ", err)
			return
		} else {
			j.B.WriteLog(string(line))
		}
	}
}
