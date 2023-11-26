package main
import (
	"net"
    "net/rpc"
    "log"
	"os"
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
	if len(os.Args) < 2 {
        log.Fatal("Usage: myapp <config_path>")
    }
    configPath := os.Args[1]

    // Initialize queue service
    queueService := NewQueueService(configPath)
    rpc.Register(queueService)
    
    // Load configuration
    config, err := LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Listen for RPC requests on the RPC address
    listener, err := net.Listen("tcp", config.RPCAddr)
    if err != nil {
        log.Fatalf("Listen error: %v", err)
    }

    log.Printf("RPC server is running on %s...\n", config.RPCAddr)
    rpc.Accept(listener)
}

