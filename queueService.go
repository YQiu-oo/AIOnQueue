package main

type QueueService struct {
    Queue *Queue
}

func NewQueueService() *QueueService {
    return &QueueService{Queue: NewQueue()}
}
type EnqueueArgs struct {
    Message string // 或者其他您需要的字段
}

type DequeueArgs struct {
    // 可以添加字段，例如 MaxMessages int
}
type EnqueueReply struct {
    Success bool // 或者其他表示操作结果的字段
}


type DequeueReply struct {
    Message string // 或者是一组消息
    Success bool   // 表示操作是否成功
}
func (q *QueueService) Enqueue(args *EnqueueArgs, reply *EnqueueReply) error {
    message := Message{Content: args.Message} 
    q.Queue.Enqueue(message)
    reply.Success = true

	return nil
    
}

func (q *QueueService) Dequeue(args *DequeueArgs, reply *DequeueReply) error {
    // 实现出队逻辑
	message := q.Queue.Dequeue()
    if message != nil {
        reply.Message = message.Content 
        reply.Success = true
    } else {
        reply.Success = false
    }
    return nil
}