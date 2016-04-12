package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/response"
	"github.com/Zombispormedio/smartdb/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func CreateSensorGrid(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}

	NewSensorGridError := sensorGrid.New(body, user, session)

	if NewSensorGridError == nil {
		response.SuccessMessage(c, "Sensor Created")
	} else {
		response.Error(c, NewSensorGridError)
	}

}