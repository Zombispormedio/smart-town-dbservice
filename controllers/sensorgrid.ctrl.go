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
		response.SuccessMessage(c, "Sensor Grid Created")
	} else {
		response.Error(c, NewSensorGridError)
	}

}

func GetSensorGrids(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	var result []models.SensorGrid

	GetAllError := models.GetSensorGrids(&result, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}

func GetSensorGridByID(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	sensorGrid := models.SensorGrid{}

	ByIdError := sensorGrid.ByID(id, session)

	if ByIdError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, ByIdError)

	}
}

func DeleteSensorGrid(c *gin.Context, session *mgo.Session) {
	id := c.Param("id")

	RemoveError := models.DeleteSensorGrid(id, session)

	if RemoveError == nil {
		GetSensorGrids(c, session)
	} else {
		response.Error(c, RemoveError)
		session.Close()
	}
}

func ChangeSensorGridSecret(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.ChangeSecret(id, session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}
}

func SetSensorGridCommunicationCenter(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")
	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.SetCommunicationCenter(id, body, session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}

}

func SetSensorGridDisplayName(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.SetDisplayName(id, body["display_name"].(string), session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}

}

func SetSensorGridZone(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.SetZone(id, body["zone"].(string), session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}

}

func AllowAccessSensorGrid(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.AllowAccess(id, session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}
}





func SetSensorGridLocation(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}

	SettingError := sensorGrid.SetLocation(id, body["location"], session)

	if SettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, SettingError)

	}

}