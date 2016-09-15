package job_test

import (
	"github.com/codestand/build/job"
	_ "github.com/codestand/build/test/testhelper"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	j := job.New()
	if !validateJobId(j.Id) {
		t.Fatal("job.Id should be uuid: ", j.Id)
	}
}

func TestSaveAndFind(t *testing.T) {
	if err := job.Save(job.Job{Id: "myjob1"}); err != nil {
		t.Fatal(err)
	}
	if err := job.Save(job.Job{Id: "myjob2"}); err != nil {
		t.Fatal(err)
	}

	if found, err := job.Find("myjob1"); err != nil {
		t.Fatal(err)
	} else if found.Id != "myjob1" {
		t.Fatal(found)
	}
}

func validateJobId(text string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	return r.MatchString(text)
}
