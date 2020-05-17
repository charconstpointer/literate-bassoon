package main

import (
	workers "alpha/api/gen"
	"alpha/domain"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	probes := make(chan domain.Probe)
	r := setupRoutes(probes)
	err := r.Run()
	if err != nil {
		log.Error("Can't start the server %v", err)
	}
}

func setupRoutes(probes chan domain.Probe) *gin.Engine {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Error("Can't connect to kafka")
	}

	go publishProbes(probes, prod)()
	r := gin.Default()
	r.POST("/probes", handleCreateProbes(probes))
	r.GET("/hello", handleHello())
	return r
}

func handleHello() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(200, "alive")
	}
}

func publishProbes(publish chan domain.Probe, prod *kafka.Producer) func() {
	return func() {
		for {
			select {
			case p := <-publish:
				invokeClient(p.Sensor)
				sendToKafka(p, prod)
			}
		}
	}
}

func sendToKafka(p domain.Probe, prod *kafka.Producer) {
	err, b := getBytes(p)
	topic := p.Sensor
	log.Infof("Publishing probe on topic %s", topic)
	if err != nil {
		log.Errorf("Can't get bytes of %v", topic)
	}
	err = prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	if err != nil {
		log.Errorf("%v could not be sent to topic", p, p.Sensor)
	}
}

func handleCreateProbes(publish chan<- domain.Probe) func(context *gin.Context) {
	return func(context *gin.Context) {
		var probes = domain.Measurement{}
		err := context.Bind(&probes)
		if err != nil {
			log.Error("Bind err")
		}
		for _, p := range probes.Probes {
			publish <- p
		}
		context.JSON(202, probes)
	}
}

func invokeClient(topic string) {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	client := workers.NewWorkersClient(cc)
	request := &workers.TopicName{Topic: topic}
	_, _ = client.ListenTopic(context.Background(), request)
}

func getBytes(p interface{}) (error, []byte) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(p)
	b := reqBodyBytes.Bytes()
	return err, b
}
