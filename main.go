package main
import "fmt"

func main() {
	fmt.Println("Hello, AIONQueue!")
	queue := NewQueue()
    producer := NewProducer()
    consumer := NewConsumer()

    // Produce some messages
    producer.Produce(queue, "Hello")
    producer.Produce(queue, "World")

    // Consume the messages
    consumer.Consume(queue)
}
