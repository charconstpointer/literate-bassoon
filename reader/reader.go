package main

import (
	"alpha/messages"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "messages",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		var probe messages.Probe
		err = json.Unmarshal(m.Value, &probe)
		if err != nil {
			fmt.Println("byte -> json err")
		}
		fmt.Printf("sensor id %d, interval %d, value %d\n", probe.SensorId, probe.Data.Interval, probe.Data.Value)
	}
}
