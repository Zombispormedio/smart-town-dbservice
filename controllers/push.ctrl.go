package controllers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/store"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
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
