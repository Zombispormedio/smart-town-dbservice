package controllers

import (
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/Zombispormedio/smartdb/models"
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

	NewZoneError := zone.New(body, user, session)

	if NewZoneError == nil {
		response.SuccessMessage(c, "Zone Created")
	} else {
		response.Error(c, NewZoneError)
	}

}

func GetZones(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	result := []models.Zone{}

	keysQueries := []string{"search", "p", "s"}
	queries := utils.Queries(c, keysQueries)

	GetAllError := models.GetZones(&result, queries, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}

}

func OtherZones(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	id := c.Param("id")
	result := []models.Zone{}

	GetAllError := models.OtherZones(&result, id, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}

}

func CountZones(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	keysQueries := []string{"search"}
	queries := utils.Queries(c, keysQueries)

	result, CountError := models.CountZones(queries, session)
	if CountError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, CountError)
	}

}

func DeleteZone(c *gin.Context, session *mgo.Session) {

	id := c.Param("id")

	RemoveError := models.DeleteZone(id, session)

	if RemoveError == nil {
		GetZones(c, session)
	} else {
		response.Error(c, RemoveError)
		session.Close()
	}

}

func GetZoneByID(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	zone := models.Zone{}

	ByIdError := zone.ByID(id, session)

	if ByIdError == nil {
		response.Success(c, zone)
	} else {
		response.Error(c, ByIdError)

	}

}

func SetZoneDisplayName(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	zone := models.Zone{}

	SettingError := zone.SetDisplayName(id, body["display_name"].(string), session)

	if SettingError == nil {
		response.Success(c, zone)
	} else {
		response.Error(c, SettingError)

	}

}

func SetZoneDescription(c *gin.Context, session *mgo.Session) {

	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	zone := models.Zone{}

	SettingError := zone.SetDescription(id, body["description"].(string), session)

	if SettingError == nil {
		response.Success(c, zone)
	} else {
		response.Error(c, SettingError)

	}
}

func SetZoneKeywords(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	zone := models.Zone{}

	SettingError := zone.SetKeywords(id, body["keywords"], session)

	if SettingError == nil {
		response.Success(c, zone)
	} else {
		response.Error(c, SettingError)

	}
}

func SetZoneShape(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	zone := models.Zone{}

	SettingError := zone.SetShape(id, body["shape"], body["center"], session)

	if SettingError == nil {
		response.Success(c, zone)
	} else {
		response.Error(c, SettingError)

	}
}

func VerifyRefZone(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	ref := c.Param("ref")

	result, RefError := models.VerifyRefZone(ref, session)
	if RefError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, RefError)
	}

}
