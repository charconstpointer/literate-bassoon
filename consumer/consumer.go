package main

import (
	workers "alpha/api/gen"
	"alpha/domain"
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	topics *map[string]bool
}

func (s *server) ListenTopic(c context.Context, t *workers.TopicName) (*workers.Empty, error) {
	topic := t.Topic
	(*s.topics)[topic] = true
	response := &workers.Empty{}
	return response, nil
}

func main() {
	_, kafkaHost, influxHost, token := parseFlags()
	probes := make(chan domain.Probe)
	topics := make(map[string]bool)
	pub := PublisherImpl{
		influxHost:  *influxHost,
		influxToken: *token,
		Probes:      probes,
	}
	reader := ReaderImpl{
		Probes:    probes,
		kafkaHost: *kafkaHost,
		topics:    &topics,
	}
	go reader.start(*kafkaHost)
	go pub.start()

	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", address)

	s := grpc.NewServer()
	workers.RegisterWorkersServer(s, &server{topics: &topics})
	s.Serve(lis)
}

func parseFlags() (*string, *string, *string, *string) {
	var topic = flag.String("topic", "D385DD88", "kafka topic")
	var kafkaHost = flag.String("kafka", "localhost:9092", "kafka host")
	var influxHost = flag.String("influx", "http://localhost:8086", "influx host")
	var token = flag.String("token", "golang:client", "influx auth")
	flag.Parse()
	log.Infof("Topic : %s, kafka : %s, influx : %s, token : %s", *topic, *kafkaHost, *influxHost, *token)
	return topic, kafkaHost, influxHost, token
}
