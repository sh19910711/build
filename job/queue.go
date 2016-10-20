package job

import (
	"github.com/codestand/build/queue"
)

var q queue.Queue
var newJobCh chan string
var finishedCh chan struct{}

func init() {
	q = queue.New("jobqueue")
	newJobCh = make(chan string)
	finishedCh = make(chan struct{})
}

func Push(buildId string) {
	q.Push(buildId)
}

func Pop() (string, error) {
	return q.Pop()
}

func QueueLength() int64 {
	return q.Len()
}

func ResetQueue() error {
	return q.Reset()
}

func Close() {
	close(finishedCh)
}