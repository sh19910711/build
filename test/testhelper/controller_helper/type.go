package controller_helper

type BuildResponse struct {
	Id  string
	Job ShowJobResponse
}

type BuildJobResponse struct {
	ExitCode int `json:exit_code`
	Finished bool
}
