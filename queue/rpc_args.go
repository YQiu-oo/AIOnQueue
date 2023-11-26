package queue

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

type QueueState struct {
    Messages []string // 或者其他你需要的信息
}