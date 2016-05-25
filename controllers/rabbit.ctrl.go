package controllers
import (
	"github.com/Zombispormedio/smartdb/models"
	"gopkg.in/mgo.v2"
    "errors"

)

func RabbitPushSensor(data map[string]interface{}, session *mgo.Session) error{
    var Error error
	
	RequestError := models.PushSensorData(session, data)

	if RequestError != nil {
		Error = errors.New(RequestError.Message)
	}

	return Error
}