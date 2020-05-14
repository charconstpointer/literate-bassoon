package main

import (
	"alpha/domain"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	probes := make(chan domain.Measurement)
	r := setupRoutes(probes)
	err := r.Run()
	if err != nil {
		log.Error("Can't start the server %v", err)
	}
}

func setupRoutes(probes chan domain.Measurement) *gin.Engine {
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

func publishProbes(publish chan domain.Measurement, prod *kafka.Producer) func() {
	return func() {
		for {
			select {
			case m := <-publish:
				sendToKafka(m, prod)
			}
		}
	}
}

func sendToKafka(m domain.Measurement, prod *kafka.Producer) {
	for _, p := range m.Probes {
		err := publishMeasurement(prod, m)
		if err != nil {
			log.Errorf("%v could not be sent to topic", p, m.Measurement)
		}
	}
}

func publishMeasurement(prod *kafka.Producer, m domain.Measurement) error {
	err, b := getBytes(m)
	topic := "measurements"
	log.Infof("Publishing probe on topic %s", topic)
	if err != nil {
		log.Errorf("Can't get bytes of %v", topic)
	}
	err = prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	return err
}

func handleCreateProbes(publish chan<- domain.Measurement) func(context *gin.Context) {
	return func(context *gin.Context) {
		var probes = domain.Measurement{}
		err := context.Bind(&probes)
		if err != nil {
			log.Error("Bind err")
		}
		publish <- probes
		context.JSON(202, probes)
	}
}

func getBytes(p interface{}) (error, []byte) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(p)
	b := reqBodyBytes.Bytes()
	return err, b
}
