package main

import (
	conn2 "alpha/conn"
	"alpha/messages"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	conn := conn2.ConnectKafka("messages", "localhost:9092", "tcp")
	listenHttp(conn)
}

func listenHttp(conn *kafka.Conn) {
	r := gin.Default()

	r.POST("/messages", handlePostMessages(conn))
	_ = r.Run()
}

func handlePostMessages(conn *kafka.Conn) func(c *gin.Context) {
	return func(c *gin.Context) {
		var probe = messages.Probe{}
		if c.Bind(&probe) == nil {
			probeBytes := getBytes(probe)
			_, err := conn.WriteMessages(kafka.Message{Value: probeBytes})
			if err != nil {
				fmt.Println("Send err")
			}
			c.JSON(200, gin.H{"body": probe.SensorId})
		}
	}
}

func getBytes(probe messages.Probe) []byte {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(probe)
	value := buffer.Bytes()
	return value
}
