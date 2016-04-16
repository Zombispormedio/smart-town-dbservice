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
		response.SuccessMessage(c, "Sensor Created")
	} else {
		response.Error(c, NewSensorError)
	}

}

func GetSensors(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")
	result:= []models.Sensor{}

	GetAllError := models.GetSensors(&result,id, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}

func GetSensorByID(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	sensor := models.Sensor{}

	ByIDError := sensor.ByID(id, session)

	if ByIDError == nil {
		response.Success(c, sensor)
	} else {
		response.Error(c, ByIDError)

	}
}


func SetSensorTransmissor(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensor := models.Sensor{}

	SettingError := sensor.SetTransmissor(id, body, session)

	if SettingError == nil {
		response.Success(c, sensor)
	} else {
		response.Error(c, SettingError)

	}
}




func SetSensorDisplayName(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensor := models.Sensor{}

	SettingError := sensor.SetDisplayName(id, body["display_name"].(string), session)

	if SettingError == nil {
		response.Success(c, sensor)
	} else {
		response.Error(c, SettingError)

	}

}


func SetSensorMagnitude(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensor := models.Sensor{}

	SettingError := sensor.SetMagnitude(id, body, session)

	if SettingError == nil {
		response.Success(c, sensor)
	} else {
		response.Error(c, SettingError)

	}

}
