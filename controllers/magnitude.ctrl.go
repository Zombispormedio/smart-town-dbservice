package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/response"
	"github.com/Zombispormedio/smartdb/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func CreateMagnitude(c *gin.Context, session *mgo.Session) {
	defer session.Close()
    
	preUser, _ := c.Get("user")
	user := preUser.(string)
    
	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMapString(bodyInterface)

	magnitude := models.Magnitude{}

	NewMagnitudeError := magnitude.New(body, user, session)

	if NewMagnitudeError == nil {
		response.SuccessMessage(c, "Magnitude Created")
	} else {
		response.Error(c, NewMagnitudeError)
	}

}

func GetMagnitudes(c *gin.Context, session *mgo.Session) {
	defer session.Close()
    
    var result []models.ListMagnitudeItem
    
    GetAllError:=models.AllMagnitudes(&result, session)
    if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}

}
