package main

import (
	"alpha/messages"
	"context"
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	"time"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "messages",
		Partition: 0,
		MinBytes:  0x3E8, // 10KB
		MaxBytes:  10e6,  // 10MB
	})
	client := influxdb2.NewClient("http://localhost:8086", "golang:client")
	// user blocking write client for writes to desired bucket

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
		writeToInflux(client, &probe)
		fmt.Printf("sensor id %d, interval %d, value %d\n", probe.SensorId, probe.Data.Interval, probe.Data.Value)
	}
}

func writeToInflux(client influxdb2.Client, probe *messages.Probe) {
	writeApi := client.WriteApiBlocking("", "probes")
	// create point using full params constructor
	p := influxdb2.NewPoint("random",
		map[string]string{"unit": "delay"},
		map[string]interface{}{"value": probe.Data.Value},
		time.Now())
	// write point immediately
	writeApi.WritePoint(context.Background(), p)
}
