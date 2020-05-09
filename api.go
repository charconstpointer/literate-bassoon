package main

import (
	"alpha/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
)

func main() {
	listenHttp()
}

func listenHttp() {
	r := gin.Default()
	r.GET("/probes", handleGetProbes("foo bar"))
	r.POST("/probes", handleCreateProbe())
	_ = r.Run()
}

func handleGetProbes(dep interface{}) func(*gin.Context) {
	fmt.Println(dep)
	return func(context *gin.Context) {
		context.JSON(200, dep)
	}
}

func handleCreateProbe() func(context *gin.Context) {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Println("Could not connect to kafka")
	}
	return func(context *gin.Context) {
		var probes = domain.Probes{}
		err := context.Bind(&probes)
		if err != nil {
			log.Panicln("Bind err")
		}
		for _, p := range probes.Probes {
			err, b := getBytes(p)
			if err != nil {
				log.Fatalf("Can't get bytes of %v", p)
			}
			err = prod.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &probes.Measurement, Partition: kafka.PartitionAny},
				Value:          b,
			}, nil)
			if err != nil {
				log.Println("Can't write to kafka")
			}
		}
		context.JSON(202, probes)
	}
}

func getBytes(p domain.Probe) (error, []byte) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(p)
	b := reqBodyBytes.Bytes()
	return err, b
}
