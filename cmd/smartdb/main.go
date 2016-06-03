package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
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

	port := os.Getenv("PORT")

	if port == "" {
		port = "5060"
	}

	if os.Getenv("IS_RABBIT") != "NO" {

		rabbit, RabbitError := consumer.New(routes.Consumer, session)

		if RabbitError != nil {
			panic(RabbitError)
		}

		routes.Set(router, session, rabbit)

		go router.Run(":" + port)

		RabbitRunError := rabbit.Run()

		if RabbitRunError != nil {
			log.Error(RabbitRunError)
		}
	} else {

		routes.Set(router, session, nil)

		router.Run(":" + port)
	}

}
