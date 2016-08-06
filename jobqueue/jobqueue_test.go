package jobqueue_test

import (
	"github.com/codestand/build/jobqueue"
	"github.com/codestand/build/test/helper"
	"testing"
)

func init() {
	helper.Init()
}

func TestPush(t *testing.T) {
	jobqueue.Init()

	newjob := jobqueue.NewJob()
	newjob.Id = "my-job"
	jobqueue.Push(newjob)

	ret := <-jobqueue.Queue()
	if ret.Id != "my-job" {
		t.Fatal(ret)
	}
}
