package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	conn2 "literate-bassoon/conn"
	"literate-bassoon/domain"
	"log"
)

func main() {
	conn := conn2.ConnectKafka("messages", "localhost:9092", "tcp")
	listenHttp(conn)
}

func listenHttp(conn *kafka.Conn) {
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
	return func(context *gin.Context) {
		var probes = domain.Probes{}
		err := context.Bind(&probes)
		if err != nil {
			log.Panicln("Bind err")
		}
		context.JSON(202, probes)
	}
}
