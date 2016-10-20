package job

import (
	"github.com/codestand/build/queue"
)

var q *queue.Queue
var finishedCh chan struct{}

func init() {
	q = queue.New("jobqueue")
	finishedCh = make(chan struct{})
}

func Push(buildId string) {
	q.Push(buildId)
}

func Pop() string {
	return q.Pop()
}

func QueueLength() int {
	return q.Len()
}

func ResetQueue() {
	q.Reset()
}

func Close() {
	close(finishedCh)
	q.Close()
}
