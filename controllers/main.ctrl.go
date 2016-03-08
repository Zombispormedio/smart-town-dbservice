package controllers

import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"

)


func Hi(c *gin.Context, session *mgo.Session){
    var msg struct{
        Message string
    }

    msg.Message="Hello World"
    c.JSON(200, msg)
}