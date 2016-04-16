package controllers

import (
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/response"
	"github.com/Zombispormedio/smartdb/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func CreateTask(c *gin.Context, session *mgo.Session) {
	defer session.Close()

	preUser, _ := c.Get("user")
	user := preUser.(string)

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	task := models.Task{}

	NewTaskError := task.New(body, user, session)

	if NewTaskError == nil {
		response.SuccessMessage(c, "Task Created")
	} else {
		response.Error(c, NewTaskError)
	}

}



func GetTasks(c *gin.Context, session *mgo.Session) {

	defer session.Close()

	result:= []models.Task{}

	GetAllError := models.GetTasks(&result, session)
	if GetAllError == nil {
		response.Success(c, result)
	} else {
		response.Error(c, GetAllError)
	}
}



func UpdateTask(c *gin.Context, session *mgo.Session) {
	defer session.Close()
	id := c.Param("id")

	bodyInterface, _ := c.Get("body")
	body := utils.InterfaceToMap(bodyInterface)

	task := models.Task{}

	SettingError := task.Update(id, body, session)

	if SettingError == nil {
		response.Success(c, task)
	} else {
		response.Error(c, SettingError)

	}
}

func DeleteTask(c *gin.Context, session *mgo.Session) {

	id := c.Param("id")

	RemoveError := models.DeleteTask(id, session)

	if RemoveError == nil {
		 GetTasks(c, session)
	} else {
		response.Error(c, RemoveError)
		session.Close()
	}

}