package queue

import (
	"github.com/adjust/redismq"
)

type Queue struct {
	client   *redismq.Queue
	consumer *redismq.Consumer
}

func New(name string) Queue {
	q := Queue{client: redismq.CreateQueue("localhost", "6379", "", 0, name)}
	q.listen()
	return q
}

func (q *Queue) Reset() error {
	for {
		if q.Len() > 0 {
			if _, err := q.Pop(); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func (q *Queue) Len() int64 {
	return q.client.GetInputLength()
}

func (q *Queue) Push(payload string) error {
	return q.client.Put(payload)
}

func (q *Queue) Pop() (msg string, err error) {
	if p, err := q.consumer.Get(); err != nil {
		return "", err
	} else {
		if err := p.Ack(); err != nil {
			return "", err
		}
		msg = p.Payload
	}
	return msg, nil
}

func (q *Queue) listen() error {
	if c, err := q.client.AddConsumer("consumer"); err != nil {
		return err
	} else {
		q.consumer = c
	}
	return nil
}
