package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/response"
    "github.com/Zombispormedio/smartdb/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)


func CreateZone(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	zone := models.Zone{}

	NewZoneError :=zone.New(body, user, session)

	if NewZoneError == nil {
		response.SuccessMessage(c, "Zone Created")
	} else {
		response.Error(c, NewZoneError)
	}

}