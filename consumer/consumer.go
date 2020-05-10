package main

import (
	"alpha/domain"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	var topic = flag.String("topic", "default", "kafka topic")
	var kafkaHost = flag.String("kafka", "localhost:9092", "kafka host")
	var influxHost = flag.String("influx", "http://localhost:8086", "influx host")
	var token = flag.String("token", "golang:client", "influx auth")
	flag.Parse()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{*kafkaHost},
		Topic:     *topic,
		Partition: 0,
		MinBytes:  0x3E8,
		MaxBytes:  10e6,
	})
	client := influxdb2.NewClient(*influxHost, *token)
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
		writeToInflux(client, &probe, *topic)
		log.Println("finished")
	}
}

func writeToInflux(client influxdb2.InfluxDBClient, probe *domain.Probe, t string) {
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
