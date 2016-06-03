package controllers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/consumer"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/store"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/Zombispormedio/smartdb/models"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
)

func PushCredentialsConfig(consumer *consumer.Consumer) func(c *gin.Context) {
	return func(c *gin.Context) {

		if consumer == nil {
			response.ErrorByString(c, 404, "No rabbit here")
		} else {

			newID, _ := uuid.NewV4()
			PushID := newID.String()

			Error := store.Put("push_identifier", PushID, "Config")

			if Error == nil {

				ResponseData := struct {
					Key string `json:"key"`
				}{
					PushID,
				}

				response.Success(c, ResponseData)

				consumer.ReBind()

			} else {
				log.WithFields(log.Fields{
					"message": Error.Error(),
				}).Error("PushCredentialsConfigError")

				response.ErrorByString(c, 404, "Error storing new key")
			}
		}
	}

}

func CheckSensorGrid(c *gin.Context, session *mgo.Session) {

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	sensorGrid := models.SensorGrid{}
	CredentialsError := sensorGrid.CheckCredentials(body["client_id"].(string), body["client_secret"].(string), session)

	if CredentialsError == nil {
		response.SuccessMessage(c, "Perfect Check in")
	} else {
		response.Error(c, CredentialsError)

	}

}

func PushSensorRegistry(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	body, _ := c.Get("body")

	Data := utils.SliceInterfaceToSliceMap(body)

	PushError := models.PushSensorData(session, Data...)

	if PushError == nil {
		response.SuccessMessage(c, "Perfect Pushover")
	} else {
		response.Error(c, PushError)

	}
}
