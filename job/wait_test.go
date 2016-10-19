package job_test

import (
	"github.com/codestand/build/job"
	"testing"
)

func TestWait(t *testing.T) {
	setup()
	defer teardown()

	getFakeBuild().ResetLog()

	// prepare jobqueue
	job.ResetQueue()
	go func() {
		job.Push("10000")
	}()
	go func() {
		sleep(1000)
		job.Close()
	}()

	job.Wait()

	// check build log
	b := getFakeBuild()
	if !contains(b.Log, "gcc -o app main.c") {
		t.Fatal("the build should run gcc command")
	}
}
