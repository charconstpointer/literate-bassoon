package main

import (
	"alpha/domain"
	"encoding/json"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type Reader interface {
	start(host string) error
}

type Publisher interface {
	start() error
}

type ReaderImpl struct {
	Probes    chan domain.Probe
	reader    *kafka.Reader
	kafkaHost string
	topics    *map[string]bool
}

type PublisherImpl struct {
	influxHost  string
	influxToken string
	Probes      chan domain.Probe
	influx      influxdb2.InfluxDBClient
	topics      *map[string]bool
}

func (r ReaderImpl) start(host string) error {
	for {
		for k, v := range *r.topics {
			fmt.Println(k, v)
			r.reader = kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{host},
				Topic:     k,
				Partition: 0,
				MinBytes:  0x3E8,
				MaxBytes:  10e6,
			})
			m, err := r.reader.ReadMessage(context.Background())
			if err != nil {
				log.Error(err)
				return err
			}
			var probe domain.Probe
			err = json.Unmarshal(m.Value, &probe)
			if err != nil {
				log.Error("Can't parse measurement")
				return err
			}
			r.Probes <- probe
			log.Println("Successfully fetched probe")
			time.Sleep(3000 * time.Millisecond)

		}
	}
}

func (c PublisherImpl) start() {
	go func() {
		for {
			select {
			case p := <-c.Probes:
				c.publishToInflux(p)
			}
		}
	}()
}

func (c PublisherImpl) publishToInflux(p domain.Probe) {
	if c.influx == nil {
		c.influx = influxdb2.NewClient(c.influxHost, c.influxToken)
	}
	writeApi := c.influx.WriteApi("", "probes")
	point := influxdb2.NewPoint(p.Sensor,
		map[string]string{"unit": "delay"},
		map[string]interface{}{p.Unit: p.Value},
		time.Now())
	writeApi.WritePoint(point)

	log.Info("Writing to influx")
}
