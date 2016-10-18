package job_test

import (
	"github.com/codestand/build/job"
	"testing"
)

func TestAttach(t *testing.T) {
	setup()
	defer teardown()

	ctx, cancel := contextWithTimeout()
	defer cancel()

	b := getFakeBuild()
	j := job.New(ctx, b)

	if err := j.Create(); err != nil {
		t.Fatal(err)
	}

	if err := j.Attach(); err != nil {
		t.Fatal(err)
	}

	if exitCode, err := j.Run(); err != nil {
		t.Fatal(err)
	} else if exitCode != 0 {
		t.Fatal("exitCode should be zero")
	}

	if !contains(j.B.Log, "gcc -o app main.c") {
		t.Fatal("the build should run gcc command")
	}
}
