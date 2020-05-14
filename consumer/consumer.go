package main

import (
	"alpha/domain"
	"alpha/messaging"
	"context"
	"encoding/json"
	"flag"
	"github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	_, kafkaHost, influxHost, token := parseFlags()
	client := influxdb2.NewClient(*influxHost, *token)
	topics := make(chan messaging.Cursor)
	done := make(chan bool)
	for i := 0; i < 2; i++ {
		go worker(*kafkaHost, client, topics, done)
	}
	topics <- messaging.Cursor{Name: "001"}
	topics <- messaging.Cursor{Name: "002"}
	topics <- messaging.Cursor{Name: "003"}
	<-done
}

func worker(h string, c influxdb2.InfluxDBClient, topics chan messaging.Cursor, done chan bool) {
	for {
		select {
		case m := <-topics:
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{h},
				Topic:     m.Name,
				Partition: 0,
				MinBytes:  0x3E8,
				MaxBytes:  10e6,
			})

			msr, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Error(err)
				time.Sleep(5000 * time.Millisecond)
			}
			var measurement domain.Measurement
			err = json.Unmarshal(msr.Value, &measurement)
			readAndPersist(measurement, c)
			go func() {
				topics <- m
				log.Info("restoring topic")
			}()
			time.Sleep(5000 * time.Millisecond)
		case _ = <-done:
			log.Info("terminating go routine")
		}
	}
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

func readAndPersist(m domain.Measurement, client influxdb2.InfluxDBClient) {
	err := writeToInflux(client, &m)
	if err != nil {
		log.Error(err)
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
