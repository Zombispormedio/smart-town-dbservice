package controllers

import(
    "github.com/gin-gonic/gin"
    "github.com/Zombispormedio/smartdb/lib/response"

)


func Hi(c *gin.Context){
    response.SuccessMessage(c, "Hello World")
}