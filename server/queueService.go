package main
import "AIONQueue/queue"
import (
	"encoding/json"
	"errors"
	// "io/ioutil"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)
type Config struct {
    RaftNodeID string `json:"RAFT_NODE_ID"`
    RaftLogPath string `json:"RAFT_LOG_PATH"`
    RaftDBPath string `json:"RAFT_DB_PATH"`
    RaftAddr string `json:"RAFT_ADDR"`
	RaftSnapshotPath string `json:"RAFT_SNAPSHOT_PATH"`
	RPCAddr string `json:"RPC_ADDR"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    
    defer file.Close()

    config := &Config{}
    decoder := json.NewDecoder(file)
    err = decoder.Decode(config)
    if err != nil {
        return nil, err
    }

    return config, nil
}




const (
    raftTimeout = 10 * time.Second
)

type QueueService struct {
    Queue *Queue
	Backup *QueueBackup
	BackupFilePath string
	Raft  *raft.Raft

}
//快照和增量结构
type QueueSnapshot struct {
    Messages []Message // 完整的队列状态
}

type QueueChange struct {
    Operation string  // "enqueue" 或 "dequeue"
    Message   Message // 涉及的消息内容
}

type QueueBackup struct {
    Snapshot *QueueSnapshot
    Changes  []QueueChange
}

func (s *QueueService) CreateSnapshot() *QueueSnapshot {
    s.Queue.lock.Lock()
    defer s.Queue.lock.Unlock()

    snapshot := &QueueSnapshot{
        Messages: make([]Message, len(s.Queue.messages)),
    }
    copy(snapshot.Messages, s.Queue.messages)
    return snapshot
}
func (s *QueueService) RecordChange(operation string, message Message) {
    s.Queue.lock.Lock()
    defer s.Queue.lock.Unlock()

    change := QueueChange{
        Operation: operation,
        Message:   message,
    }
    // 假设 s.Backup 是 QueueBackup 类型的实例
    s.Backup.Changes = append(s.Backup.Changes, change)
}



func checkIfFirstLaunch(dbPath string) bool {
    // Check if the Raft database file exists
    if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
        // File does not exist, indicating a first launch
        log.Println("First launch detected, no existing Raft DB found.")
        return true
    }
    // File exists, not a first launch
    log.Println("Raft DB exists, not a first launch.")
    return false
}

func NewQueueService(configPath string) *QueueService {
	config, err := LoadConfig(configPath)
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }


    service := &QueueService{
        Queue: NewQueue(),
        Backup: &QueueBackup{Snapshot: nil, Changes: []QueueChange{}},
        BackupFilePath: configPath,
    }

    fsm := &QueueFSM{QueueService: service}

    // 从环境变量中读取配置
    raftAddr := config.RaftAddr
    raftNodeID := config.RaftNodeID
    raftLogPath := config.RaftLogPath
    raftDBPath := config.RaftDBPath

    // 创建并配置 Raft 实例
    cfg := raft.DefaultConfig()
    cfg.LocalID = raft.ServerID(raftNodeID)

	isFirstLaunch := checkIfFirstLaunch(config.RaftSnapshotPath)

    // 日志存储
    logStore, err := raftboltdb.NewBoltStore(raftLogPath)
    if err != nil {
        log.Fatalf("failed to create log store: %v", err)
    }

    // 稳定存储
    stableStore, err := raftboltdb.NewBoltStore(raftDBPath)
    if err != nil {
        log.Fatalf("failed to create stable store: %v", err)
    }

	snapshotStore, err := raft.NewFileSnapshotStore(config.RaftSnapshotPath, 1, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create snapshot store: %v", err)
	}

    // 传输
    addr, err := net.ResolveTCPAddr("tcp", raftAddr)
    if err != nil {
        log.Fatalf("failed to resolve TCP address: %v", err)
    }
    transport, err := raft.NewTCPTransport(addr.String(), addr, 3, raftTimeout, os.Stderr)
    if err != nil {
        log.Fatalf("failed to create transport: %v", err)
    }

    // 初始化 Raft 实例
    raftInstance, err := raft.NewRaft(cfg, fsm, logStore, stableStore, snapshotStore, transport)
    if err != nil {
        log.Fatalf("failed to create raft: %v", err)
    }
	
    if isFirstLaunch {
        configuration := raft.Configuration{
            Servers: []raft.Server{
                {
                    ID:      raft.ServerID("node1"),
                    Address: raft.ServerAddress("localhost:1001"),
                },
                {
                    ID:      raft.ServerID("node2"),
                    Address: raft.ServerAddress("localhost:1002"),
                },
                {
                    ID:      raft.ServerID("node3"),
                    Address: raft.ServerAddress("localhost:1003"),
                },
            },
        }
        raftInstance.BootstrapCluster(configuration)
    }

    service.Raft = raftInstance

	return service


}



func (q *QueueService) Enqueue(args *queue.EnqueueArgs, reply *queue.EnqueueReply) error {

	if args.Message == "" {
        return errors.New("message cannot be empty")
    }
    message := Message{Content: args.Message} 
	// q.RecordChange("enqueue", message)
    q.Queue.Enqueue(message)
	op := QueueChange{Operation: "enqueue", Message: message}
	data, err := json.Marshal(op)
	if err != nil {
        return err
    }
    reply.Success = true
	q.Raft.Apply(data, 500*time.Millisecond)

	return nil
    
}

func (q *QueueService) Dequeue(args *queue.DequeueArgs, reply *queue.DequeueReply) error {
    // 实现出队逻辑
	message := q.Queue.Dequeue()
	m := Message{Content: message.Content}
	op := QueueChange{Operation: "dequeue", Message: m}
    data, err := json.Marshal(op)
    if err != nil {
        return err
    }
	q.Raft.Apply(data, 500*time.Millisecond) // 调整超时时间
    if message != nil {
        reply.Message = message.Content 
        reply.Success = true
	// 	q.RecordChange("dequeue", Message{Content: reply.Message})
    // } else {
        reply.Success = false
    }
    return nil
}


// func (s *QueueService) SaveBackupToFile(filePath string) error {
//     //序列化并保存 s.Backup 到文件
//     s.Queue.lock.Lock()
//     defer s.Queue.lock.Unlock()

//     data, err := json.Marshal(s.Backup)
//     if err != nil {
//         return err
//     }

//     return ioutil.WriteFile(filePath, data, 0644)
// }


// func (s *QueueService) LoadBackupFromFile(filePath string) error {
//     //从文件加载并反序列化到 s.Backup

// 	s.Queue.lock.Lock()
//     defer s.Queue.lock.Unlock()

//     data, err := ioutil.ReadFile(filePath)
//     if err != nil {
//         return err
//     }

//     err = json.Unmarshal(data, &s.Backup)
//     if err != nil {
//         return err
//     }
//     return nil
    
// }


// func (s *QueueService) StartSnapshotRoutine(interval time.Duration) {
//     ticker := time.NewTicker(interval)
//     for {
//         select {
//         case <-ticker.C:
//             s.Backup.Snapshot = s.CreateSnapshot()
//             s.Backup.Changes = []QueueChange{}
//             err := s.SaveBackupToFile(s.BackupFilePath)
//             if err != nil {
//                 log.Printf("Error saving backup: %v", err)
//             }
//         }
//     }
// }

type QueueFSM struct {
    QueueService *QueueService
}


func (f *QueueFSM) Apply(lg *raft.Log) interface{} {
    var op QueueChange
    if err := json.Unmarshal(lg.Data, &op); err != nil {
        log.Printf("Failed to unmarshal log data: %v", err)
        return nil
    }

    switch op.Operation {
    case "enqueue":
        f.QueueService.Queue.Enqueue(Message{Content: op.Message.Content})
    case "dequeue":
        f.QueueService.Queue.Dequeue()
    }
    return nil
}

func (f *QueueFSM) Restore(rc io.ReadCloser) error {
    // 从 ReadCloser 中读取并恢复快照
    var snapshot QueueSnapshot
    if err := json.NewDecoder(rc).Decode(&snapshot); err != nil {
        return err
    }

    f.QueueService.Queue.messages = snapshot.Messages
    return nil
}

type QueueSnapshotter struct {
    Snap *QueueSnapshot
}


func (s *QueueSnapshotter) Persist(sink raft.SnapshotSink) error {
    err := func() error {
        // Encode the snapshot and write to the sink
        if err := json.NewEncoder(sink).Encode(s.Snap); err != nil {
            return err
        }
        return sink.Close()
    }()

    if err != nil {
        sink.Cancel() // Cancel the snapshot operation in case of error
    }
    return err
}

func (s *QueueSnapshotter) Release() {
    // Release resources or perform any cleanup if needed
}
func (f *QueueFSM) Snapshot() (raft.FSMSnapshot, error) {
    // 从 QueueService 获取当前队列的快照
    snapshot := f.QueueService.CreateSnapshot()
    return &QueueSnapshotter{Snap: snapshot}, nil
}




func (q *QueueService) GetQueueState(args *struct{}, reply *queue.QueueState) error {
    q.Queue.lock.Lock()
    defer q.Queue.lock.Unlock()

    for _, msg := range q.Queue.messages {
        reply.Messages = append(reply.Messages, msg.Content)
    }
    return nil
}           