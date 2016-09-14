package job_test

import (
	"github.com/codestand/build/job"
	"regexp"
	"testing"
)

func TestNewJob(t *testing.T) {
	j := job.New()
	if !validateId(j.Id) {
		t.Fatal("job.Id should be uuid: ", j.Id)
	}
}

func validateId(text string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	return r.MatchString(text)
}
