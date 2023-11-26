package main

import (
    "testing"
    "net/rpc"
    "log"
    "AIONQueue/queue"
)

func TestQueueService(t *testing.T) {
    // 连接到 RPC 服务器
    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("Connection Errors", err)
    }

    // 测试 Enqueue 方法
    enqueueArgs := queue.EnqueueArgs{Message: "Hello, AIONQueue!"}
    var enqueueReply queue.EnqueueReply
    err = client.Call("QueueService.Enqueue", &enqueueArgs, &enqueueReply)
    if err != nil || !enqueueReply.Success {
        log.Fatal("Enqueue fails:", err)
    }
    log.Println("Enqueue successes")

    // 测试 Dequeue 方法
    dequeueArgs := queue.DequeueArgs{}
    var dequeueReply queue.DequeueReply
    err = client.Call("QueueService.Dequeue", &dequeueArgs, &dequeueReply)
    if err != nil || !dequeueReply.Success {
        log.Fatal("Dequeue falis", err)
    }
    log.Printf("Dequeue successes%s\n", dequeueReply.Message)
}