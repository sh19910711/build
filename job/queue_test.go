package job_test

import (
	"github.com/codestand/build/job"
	"testing"
)

func TestPushAndPop(t *testing.T) {
	job.ResetQueue()

	job.Push("this-is-build-id")
	job.Push("this-is-another-build-id")

	if buildId, err := job.Pop(); err != nil {
		t.Fatal(err)
	} else if buildId != "this-is-build-id" {
		t.Fatal("the-build-id is wrong: " + buildId)
	}

	if buildId, err := job.Pop(); err != nil {
		t.Fatal(err)
	} else if buildId != "this-is-another-build-id" {
		t.Fatal("the-build-id is wrong: " + buildId)
	}
}
