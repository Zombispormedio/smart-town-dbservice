package controllers

import (
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func Hi(c *gin.Context) {
	response.SuccessMessage(c, "Hello World")
}

func Status(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	MongoError := session.Ping()

	status := struct {
		DB bool `json:"db_status"`
	}{}

	status.DB = MongoError == nil

	response.Success(c, status)

}
