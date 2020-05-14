package main

import (
	"alpha/domain"
	"context"
	"encoding/json"
	"flag"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	topic, kafkaHost, influxHost, token := parseFlags()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{*kafkaHost},
		Topic:     *topic,
		Partition: 0,
		MinBytes:  0x3E8,
		MaxBytes:  10e6,
	})
	client := influxdb2.NewClient(*influxHost, *token)
	readAndPersist(r, client)
}

func parseFlags() (*string, *string, *string, *string) {
	var topic = flag.String("topic", "measurements", "kafka topic")
	var kafkaHost = flag.String("kafka", "localhost:9092", "kafka host")
	var influxHost = flag.String("influx", "http://localhost:8086", "influx host")
	var token = flag.String("token", "golang:client", "influx auth")
	flag.Parse()
	log.Infof("Topic : %s, kafka : %s, influx : %s, token : %s", *topic, *kafkaHost, *influxHost, *token)
	return topic, kafkaHost, influxHost, token
}

func readAndPersist(r *kafka.Reader, client influxdb2.InfluxDBClient) {
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Error(err)
			time.Sleep(5000 * time.Millisecond)
		}
		var measurement domain.Measurement
		err = json.Unmarshal(m.Value, &measurement)
		if err != nil {
			log.Error("Can't parse measurement")
		} else {
			err = writeToInflux(client, &measurement)
			if err != nil {
				log.Error(err)
			}
		}

	}
}
func writeToInflux(client influxdb2.InfluxDBClient, m *domain.Measurement) error {
	topic := m.Measurement
	writeApi := client.WriteApi("", "probes")
	for _, p := range m.Probes {
		point := influxdb2.NewPoint(topic,
			map[string]string{"unit": "delay"},
			p.Values,
			time.Now())
		writeApi.WritePoint(point)
	}
	log.Info("Writing to influx")
	return nil
}
