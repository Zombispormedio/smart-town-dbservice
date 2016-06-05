package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/utils"
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

	keysQueries := []string{"search", "p", "s"}
	queries := utils.Queries(c, keysQueries)

	GetAllError := models.GetSensors(&result,id, queries, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}

func CountSensors(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	defer session.Close()
	id := c.Param("id")
	
	keysQueries:=[]string{"search"}
	queries:=utils.Queries(c, keysQueries)

	result, CountError := models.CountSensors(id, queries, session)
	if CountError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, CountError)
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


func ImportSensors(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	
	body:=utils.SliceInterfaceToSliceMap(bodyInterface)
	
	ImportError := models.ImportSensors(body, user, session)
	
	if ImportError == nil {
		response.SuccessMessage(c, "Sensors Imported")
	} else {
		response.Error(c, ImportError)
	
	}
}

func SensorsNotifications(c *gin.Context, session *mgo.Session){
	defer session.Close()
	
	result:= []models.Sensor{}
	
	GetAllError := models.SensorsNotifications(&result, session)
	
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}

func FixSensor(c *gin.Context, session *mgo.Session){
	id := c.Param("id")

	FixError := models.FixSensor(id, session)

	if FixError == nil {
		SensorsNotifications(c, session)
	} else {
		response.Error(c, FixError)
		session.Close()
	}
}

func ReviewLastSync(c *gin.Context, session *mgo.Session){
	defer session.Close()
	
	
	ReviewError := models.ReviewSensorLastSync(session)

	if ReviewError == nil {
		response.SuccessMessage(c, "Review Successful")
	} else {
		response.Error(c, ReviewError)
		
	}
}