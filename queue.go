package main

import "sync"


type Message struct {
	Content string
}

type Queue struct {
	messages []Message
	lock sync.Mutex 
}

func NewQueue() *Queue {
	return &Queue {
		messages: []Message{},
	}
}

func (q *Queue) Enqueue(message Message) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.messages = append(q.messages, message)

}

func (q *Queue) Dequeue() *Message {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.messages) > 0 {
		message := q.messages[0]
		q.messages = q.messages[1:]
		return &message
	}
	return nil
}