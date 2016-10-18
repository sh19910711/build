package job

func (j *Job) Run() (exitCode int, err error) {
	if err := j.w.Start(j.ctx); err != nil {
		return -1, err
	}

	if exitCode, err = j.w.Wait(j.ctx); err != nil {
		return -2, err
	}

	return exitCode, nil
}
