package main

import (
	"alpha/domain"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	"time"
)

func main() {
	var topic = flag.String("topic", "t4", "kafka topic")
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
			fmt.Println(err)
			time.Sleep(5000 * time.Millisecond)
		}
		var probe domain.Probe
		err = json.Unmarshal(m.Value, &probe)
		if err != nil {
			fmt.Println("cant parse probe")
		} else {
			err = writeToInflux(client, &probe, *topic)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func writeToInflux(client influxdb2.InfluxDBClient, probe *domain.Probe, t string) error {
	writeApi := client.WriteApiBlocking("", "probes")
	p := influxdb2.NewPoint(t,
		map[string]string{"unit": "delay"},
		map[string]interface{}{"value": probe.Value},
		time.Now())
	err := writeApi.WritePoint(context.Background(), p)
	return err
}
