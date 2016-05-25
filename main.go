package main

import (
	"os"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/consumer"
	"github.com/Zombispormedio/smartdb/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	config.ServerConfig(router)

	session := config.SessionDB()

	if session == nil {
		panic("MongodbSession Fault")
	} else {
		defer session.Close()
	}

	rabbit, RabbitError := consumer.New()

	if RabbitError != nil {
		panic(RabbitError)
	}

	routes.Set(router, session, rabbit)

	port := os.Getenv("PORT")

	if port == "" {
		port = "5060"
	}

	go router.Run(":" + port)
	rabbit.Run()
}
