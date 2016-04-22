package controllers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/store"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/Zombispormedio/smartdb/models"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
)

func PushCredentialsConfig(c *gin.Context) {

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

	} else {
		log.WithFields(log.Fields{
			"message": Error.Error(),
		}).Error("PushCredentialsConfigError")

		response.ErrorByString(c, 404, "Error storing new key")
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
		response.Error(c,CredentialsError)

	}
	

}
