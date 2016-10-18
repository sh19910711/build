// +build integration

package job_test

import (
	"github.com/codestand/build/job"
	"testing"
)

func TestRun(t *testing.T) {
	setup()
	defer teardown()

	ctx, cancel := contextWithTimeout()
	defer cancel()

	b := getFakeBuild()
	j := job.New(ctx, b)

	if err := j.Create(); err != nil {
		t.Fatal(err)
	}

	if exitCode, err := j.Run(); err != nil {
		t.Fatal(err)
	} else if exitCode != 0 {
		t.Fatal("exitCode should be zero")
	}
}
