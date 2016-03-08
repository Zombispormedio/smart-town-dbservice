package controllers
import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "github.com/Zombispormedio/smartdb/response"

)

func Register(c *gin.Context, session *mgo.Session ){

    defer session.Close()

    response.SuccessMessage(c, "Hello World")

}