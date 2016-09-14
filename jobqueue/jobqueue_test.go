package jobqueue_test

import (
	"github.com/codestand/build/job"
	"github.com/codestand/build/jobqueue"
	_ "github.com/codestand/build/test/testhelper"
	"testing"
)

func TestPush(t *testing.T) {
	go func() {
		jobqueue.Push(job.Job{Id: "job1"})
		jobqueue.Push(job.Job{Id: "job2"})
		jobqueue.Push(job.Job{Id: "job3"})
	}()

	q := jobqueue.Queue()
	defer jobqueue.Close()
	if ret := (<-q).Id; ret != "job1" {
		t.Fatal(ret)
	}
	if ret := (<-q).Id; ret != "job2" {
		t.Fatal(ret)
	}
	if ret := (<-q).Id; ret != "job3" {
		t.Fatal(ret)
	}
}
