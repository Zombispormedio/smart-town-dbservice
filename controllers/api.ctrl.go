package controllers
import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "github.com/Zombispormedio/smartdb/response"
    "github.com/Zombispormedio/smartdb/models"

)


func Register(c *gin.Context, session *mgo.Session ){

    defer session.Close()

    var body map[string]string

    oauth:=models.OAuth{}
    if c.BindJSON(&body)==nil{



        oauth.New(body, session)  
        response.SuccessMessage(c, "all good")



    }
}