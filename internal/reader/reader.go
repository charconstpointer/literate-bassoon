package reader

import (
	"alpha/domain"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type Reader interface {
	Start(host string) error
}

type ReaderImpl struct {
	Probes    chan domain.Probe
	reader    *kafka.Reader
	KafkaHost string
	Topics    *map[string]bool
}

func (r ReaderImpl) Start(host string) error {
	for {
		for k, v := range *r.Topics {
			if v {
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
}
