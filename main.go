package main
import (
	"net"
    "net/rpc"
    "log"
)

func main() {
	// fmt.Println("Hello, AIONQueue!")
	// queue := NewQueue()
    // producer := NewProducer()
    // consumer := NewConsumer()

    // // Produce some messages
    // producer.Produce(queue, "Hello")
    // producer.Produce(queue, "World")

    // // Consume the messages
    // consumer.Consume(queue)
    queueService := NewQueueService()
    rpc.Register(queueService)

	
	
    listener, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("Listen fails", err)
    }

    log.Println("RPC server is running...")
    rpc.Accept(listener)
}
