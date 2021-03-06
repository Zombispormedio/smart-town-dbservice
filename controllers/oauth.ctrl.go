package controllers

import (
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/Zombispormedio/smartdb/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func Register(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	oauth := models.OAuth{}

	NewOauthError := oauth.Register(body, session)

	if NewOauthError == nil {
		response.SuccessMessage(c, "User Registered")
	} else {
		response.Error(c, NewOauthError)
	}

}

func Login(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	token, LoginError := models.Login(body, session)

	if LoginError == nil {
		response.Success(c, token)
	} else {
		response.Error(c, LoginError)
	}

}

func Logout(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	token := c.Request.Header.Get("Authorization")
	preUser, _ := c.Get("user")
	user := preUser.(string)

	LogoutError := models.Logout(token, user, session)

	if LogoutError == nil {
		response.SuccessMessage(c, "Congratulations, You have logged out")
	} else {
		response.Error(c, LogoutError)
	}
}

func Whoiam(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)
	var profile = models.Profile{}

	WhoiamError := models.GetProfile(user, &profile, session)

	if WhoiamError == nil {
		response.Success(c, profile)
	} else {
		response.Error(c, WhoiamError)
	}
}

func SetOauthDisplayName(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	preUser, _ := c.Get("user")
	id := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	var profile = models.Profile{}
	SettingError := profile.SetDisplayName(id, body["display_name"].(string), session)

	if SettingError == nil {
		response.Success(c, profile)
	} else {
		response.Error(c, SettingError)

	}

}

func SetOauthEmail(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	preUser, _ := c.Get("user")
	id := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	var profile = models.Profile{}
	SettingError := profile.SetEmail(id, body["email"].(string), session)

	if SettingError == nil {
		response.Success(c, profile)
	} else {
		response.Error(c, SettingError)

	}

}

func SetOauthPassword(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	preUser, _ := c.Get("user")
	id := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	var profile = models.Profile{}
	SettingError := profile.SetPassword(id, body["password"].(string), session)

	if SettingError == nil {
		response.Success(c, profile)
	} else {
		response.Error(c, SettingError)

	}

}

func DeleteOauth(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	id := preUser.(string)

	RemoveError := models.DeleteOauth(id, session)

	if RemoveError == nil {
		response.SuccessMessage(c, "Congratulations, You are out :'(")
	} else {
		response.Error(c, RemoveError)

	}
}

func Invite(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)
	email := body["email"].(string)
	code, InvitationError := models.Invite(email, session)

	if InvitationError != nil {
		response.Error(c, InvitationError)
		return
	}
	SendingError := utils.SendInvitation(code, email)

	if SendingError == nil {
		response.SuccessMessage(c, "Congratulations, Your invitation have sent!")
	} else {
		response.Error(c, SendingError)
	}

}

func CheckInvitation(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	code := c.Param("code")

	checked, CheckingError := models.CheckInvitation(code, session)

	if CheckingError == nil {
		response.Success(c, map[string]bool{
			"success": checked,
		})
	} else {
		response.Error(c, CheckingError)

	}

}

func Invitation(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	code := c.Param("code")
	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)
	password := body["password"].(string)

	AcceptedError := models.AcceptInvitation(code, password, session)

	if AcceptedError == nil {
		response.SuccessMessage(c, "Congratulations, Your are in!")
	} else {
		response.Error(c, AcceptedError)

	}

}
