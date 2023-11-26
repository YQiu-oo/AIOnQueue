package main

import "fmt"

type Producer struct {
	//necessary fields
}

// NewProducer
func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Produce(queue *Queue, messageContent string) {
	//Implement message sending logic
	message := Message{Content: messageContent}
	queue.Enqueue(message)
	fmt.Println("Produced:", messageContent)
}
