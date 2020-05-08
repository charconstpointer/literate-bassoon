package conn

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func ConnectKafka(topic string, host string, protocol string) *kafka.Conn {
	partition := 0
	conn, _ := kafka.DialLeader(context.Background(), protocol, host, topic, partition)
	return conn
}
