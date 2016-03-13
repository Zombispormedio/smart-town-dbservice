package controllers
import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "github.com/Zombispormedio/smartdb/response"
    "github.com/Zombispormedio/smartdb/models"
    "github.com/Zombispormedio/smartdb/utils"
    // "fmt"
)



func Register(c *gin.Context, session *mgo.Session ){

    defer session.Close()

    body_interface,_:=c.Get("body")
    body:=utils.InterfaceToMapString(body_interface);

    oauth:=models.OAuth{}    

    NewOauthError:=oauth.Register(body, session)

    if  NewOauthError == nil{
        response.SuccessMessage(c, "User Registered")
    }else{
        response.Error(c, NewOauthError);
    }

}


func Login(c *gin.Context, session *mgo.Session ){
    defer session.Close()

    body_interface,_:=c.Get("body")
    body:=utils.InterfaceToMapString(body_interface);
    
    token, LoginError:=models.Login(body, session)
    
    if LoginError == nil{
        response.Success(c, token)
    }else{
           response.Error(c, LoginError);
    }
    

}

func Logout(c *gin.Context, session *mgo.Session){
     defer session.Close()
     response.SuccessMessage(c, "Perfect")
}