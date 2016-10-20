package queue

// TODO: use redis
type Queue struct {
	queue chan string
}

func New(name string) *Queue {
	return &Queue{queue: make(chan string, 10)}
}

func (q *Queue) Reset() {
	for len(q.queue) > 0 {
		<-q.queue
	}
}

func (q *Queue) Len() int {
	return len(q.queue)
}

func (q *Queue) Push(item string) {
	q.queue <- item
}

func (q *Queue) Pop() string {
	return <-q.queue
}

func (q *Queue) Close() {
	close(q.queue)
}
