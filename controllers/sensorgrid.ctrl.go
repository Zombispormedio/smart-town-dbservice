package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/utils"
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

	result:= []models.SensorGrid{}
	
	
	keysQueries := []string{"search", "p", "s"}
	queries := utils.Queries(c, keysQueries)

	GetAllError := models.GetSensorGrids(&result, queries, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}


func CountSensorGrids(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	keysQueries:=[]string{"search"}
	queries:=utils.Queries(c, keysQueries)

	result, CountError := models.CountSensorGrids(queries, session)
	if CountError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, CountError)
	}
	
}

func GetSensorGridByID(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	sensorGrid := models.SensorGrid{}

	ByIDError := sensorGrid.ByID(id, session)

	if ByIDError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, ByIDError)

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


func DeleteSensorByGrid(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")
	sensorID := c.Param("sensor_id")
	
	sensorGrid := models.SensorGrid{}

	UnSettingError := sensorGrid.UnsetSensor(id,sensorID, session)

	if UnSettingError == nil {
		response.Success(c, sensorGrid)
	} else {
		response.Error(c, UnSettingError)

	}
}





