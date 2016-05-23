package controllers

import(
    "github.com/gin-gonic/gin"
    "github.com/Zombispormedio/smartdb/lib/response"
"gopkg.in/mgo.v2"
)


func Hi(c *gin.Context){
    response.SuccessMessage(c, "Hello World")
}

func Status(c *gin.Context, session *mgo.Session) {

	defer session.Close()
    

	MongoError:=session.Ping()
    
    status:=struct{
        DB bool `json:"db_status"`
    }{}
    
    status.DB=MongoError==nil
    
    response.Success(c, status)
    
}