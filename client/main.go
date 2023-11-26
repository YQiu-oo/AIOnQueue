package main

import (
    "net/rpc"
    "log"
    "AIONQueue/queue"


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

    // 对领导者节点发送 RPC 请求
    sendRPCRequest("localhost:1234", " nb   World")

    // 验证其他节点的状态
    // 这可能需要你实现一个新的 RPC 方法来获取队列状态
    verifyNodeState("localhost:1235")
    verifyNodeState("localhost:1236")
}







func sendRPCRequest(address string, message string) {
    client, err := rpc.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Dialing error: %v", err)
    }
    defer client.Close()

    args := queue.EnqueueArgs{Message: message}
    var reply queue.EnqueueReply
    err = client.Call("QueueService.Enqueue", args, &reply)
    if err != nil {
        log.Fatalf("QueueService error: %v", err)
    }
    log.Printf("Enqueue response from %s: %v", address, reply)
}

func verifyNodeState(address string) {
    client, err := rpc.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Dialing error: %v", err)
    }
    defer client.Close()

    var state queue.QueueState
    err = client.Call("QueueService.GetQueueState", new(struct{}), &state)
    if err != nil {
        log.Fatalf("QueueService error: %v", err)
    }

    // 这里你可以添加逻辑来验证队列的状态是否符合预期
    // 例如，检查 state.Messages 是否包含特定的消息
    log.Printf("Queue state from %s: %+v", address, state)
}