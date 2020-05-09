package main

import (
	"alpha/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	topic := "t1"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0,
		MinBytes:  0x3E8,
		MaxBytes:  10e6,
	})
	client := influxdb2.NewClient("http://localhost:8086", "golang:client")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		var probe domain.Probe
		err = json.Unmarshal(m.Value, &probe)
		if err != nil {
			fmt.Println("byte -> json err")
		}
		writeToInflux(client, &probe, topic)
		log.Println("finished")
	}
}

func writeToInflux(client influxdb2.Client, probe *domain.Probe, t string) {
	writeApi := client.WriteApiBlocking("", "probes")
	p := influxdb2.NewPoint(t,
		map[string]string{"unit": "delay"},
		map[string]interface{}{"value": probe.Value},
		time.Now())
	err := writeApi.WritePoint(context.Background(), p)
	if err != nil {
		log.Println("Could not write to influx")
	}
}
