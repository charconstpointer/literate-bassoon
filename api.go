package main

import (
	"alpha/domain"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
)

func main() {
	r := setupRoutes()
	err := r.Run()
	if err != nil {
		log.Fatalf("Can't start the server %v", err)
	}
}

func setupRoutes() *gin.Engine {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Can't connect to kafka, %v", err)
	}
	publish := make(chan domain.Measurement)
	go publishProbes(publish, prod)()
	r := gin.Default()
	r.POST("/probes", handleCreateProbe(publish))
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(200, "alive")
	})
	return r
}

func publishProbes(publish chan domain.Measurement, prod *kafka.Producer) func() {
	return func() {
		for {
			select {
			case m := <-publish:
				for _, p := range m.Probes {
					err := publishProbe(p, prod, m.Measurement)
					if err != nil {
						log.Fatalf("%v could not be sent to topic %s \n", p, m.Measurement)
					}
				}
			}
		}
	}
}

func handleCreateProbe(publish chan<- domain.Measurement) func(context *gin.Context) {
	return func(context *gin.Context) {
		var probes = domain.Measurement{}
		err := context.Bind(&probes)
		if err != nil {
			log.Panicln("Bind err")
		}
		publish <- probes
		context.JSON(202, probes)
	}
}

func publishProbe(p domain.Probe, prod *kafka.Producer, topic string) error {
	err, b := getBytes(p)
	if err != nil {
		log.Fatalf("Can't get bytes of %v", p)
	}
	err = prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	return err
}

func getBytes(p domain.Probe) (error, []byte) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(p)
	b := reqBodyBytes.Bytes()
	return err, b
}
