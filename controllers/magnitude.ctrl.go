package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func CreateMagnitude(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

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

	result:= []models.Magnitude{}
	
	keysQueries:=[]string{"search", "p", "s"}
	queries:=utils.Queries(c, keysQueries)

	GetAllError := models.GetMagnitudes(&result, queries, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}

}

func CountMagnitudes(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	keysQueries:=[]string{"search"}
	queries:=utils.Queries(c, keysQueries)

	result, CountError := models.CountMagnitudes(queries, session)
	if CountError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, CountError)
	}
	
}

func DeleteMagnitude(c *gin.Context, session *mgo.Session) {
	id := c.Param("id")

	RemoveError := models.DeleteMagnitude(id, session)

	if RemoveError == nil {
		GetMagnitudes(c, session)
	} else {
		response.Error(c, RemoveError)
		session.Close()
	}

}

func GetMagnitudeByID(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	magnitude := models.Magnitude{}

	ByIdError := magnitude.ByID(id, session)

	if ByIdError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, ByIdError)

	}

}

func SetMagnitudeDisplayName(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	magnitude := models.Magnitude{}

	SettingError := magnitude.SetDisplayName(id, body["display_name"].(string), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func SetMagnitudeType(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")

	body := utils.InterfaceToMap(bodyInterface)

	magnitude := models.Magnitude{}

	SettingError := magnitude.SetType(id, body["type"].(string), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func SetMagnitudeDigitalUnits(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	magnitude := models.Magnitude{}

	SettingError := magnitude.SetDigitalUnits(id, body["digital_units"].(map[string]interface{}), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func AddMagnitudeAnalogUnit(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	magnitude := models.Magnitude{}

	SettingError := magnitude.AddAnalogUnit(id, body["analog_unit"].(map[string]interface{}), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func UpdateMagnitudeAnalogUnit(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")
	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)
	magnitude := models.Magnitude{}

	SettingError := magnitude.UpdateAnalogUnit(id, body["analog_unit"].(map[string]interface{}), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func DeleteMagnitudeAnalogUnit(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")
	analogID := c.Param("analog_id")

	magnitude := models.Magnitude{}

	RemoveError := magnitude.DeleteAnalogUnit(id, analogID, session)

	if RemoveError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, RemoveError)

	}

}


func AddMagnitudeConversion(c *gin.Context, session *mgo.Session){
    defer session.Close()
	id := c.Param("id")
    
    bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	magnitude := models.Magnitude{}

	SettingError := magnitude.AddConversion(id, body["conversion"].(map[string]interface{}), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}


func UpdateMagnitudeConversion(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")
	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)
	magnitude := models.Magnitude{}

	SettingError := magnitude.UpdateConversion(id, body["conversion"].(map[string]interface{}), session)

	if SettingError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, SettingError)

	}
}

func DeleteMagnitudeConversion(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")
	conversionID := c.Param("conversion_id")

	magnitude := models.Magnitude{}

	RemoveError := magnitude.DeleteConversion(id, conversionID, session)

	if RemoveError == nil {
		response.Success(c, magnitude)
	} else {
		response.Error(c, RemoveError)

	}

}


func VerifyRefMagnitude(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	ref := c.Param("ref")

	result, RefError := models.VerifyRefMagnitude(ref, session)
	if RefError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, RefError)
	}
	
}