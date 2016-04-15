package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/response"
	"github.com/Zombispormedio/smartdb/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func CreateSensor(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensor := models.Sensor{}

	NewSensorError := sensor.New(body, user, session)

	if NewSensorError == nil {
		response.SuccessMessage(c, "Sensorr Created")
	} else {
		response.Error(c, NewSensorError)
	}

}