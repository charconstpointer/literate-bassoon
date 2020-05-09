package main

import (
	"alpha/commands"
	conn2 "alpha/conn"
	"alpha/domain"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	conn := conn2.ConnectKafka("messages", "localhost:9092", "tcp")
	listenHttp(conn)
}

func listenHttp(conn *kafka.Conn) {
	r := gin.Default()
	r.POST("/probes", handleCreateProbe())
	_ = r.Run()
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

func parseUser(context *gin.Context) commands.CreateUser {
	var user = commands.CreateUser{}
	err := context.Bind(&user)
	if err != nil {
		log.Panicln("cant parse user")
	}
	return user
}
