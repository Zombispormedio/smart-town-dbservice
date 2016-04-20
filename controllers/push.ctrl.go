package controllers

import (
	"github.com/Zombispormedio/smartdb/response"
	"github.com/gin-gonic/gin"
)

func PushCredentialsConfig(c *gin.Context) {
	response.SuccessMessage(c, "Hello World")
}
