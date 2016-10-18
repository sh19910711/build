package queue_test

import (
	"github.com/codestand/build/queue"
	"testing"
)

func TestPush(t *testing.T) {
	q := queue.New("test-push")
	q.Push("hello1")
	q.Push("hello2")
	q.Push("hello3")

	for _, expected := range []string{"hello1", "hello2", "hello3"} {
		if msg, err := q.Pop(); err != nil {
			t.Fatal(err)
		} else if msg != expected {
			t.Fatal("msg should be hello")
		}
	}
}
