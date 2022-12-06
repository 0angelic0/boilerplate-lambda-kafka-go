package kafka

import "time"

// Message struct is represented a message in event bus
type Message struct {
	Topic     string
	Partition int64
	Offset    int64
	Headers   []map[string][]byte
	Key       []byte
	Value     []byte
	Timestamp time.Time
}
