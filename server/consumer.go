package main

import "fmt"

type Consumer struct {
	//necessary fields
}

// NewConsumer
func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Consume(queue *Queue) {
	//Implement message receiving  logic
	for {
		message := queue.Dequeue()
		if message != nil {
			fmt.Println("Consumed: ", message.Content)
		}
	}
}
