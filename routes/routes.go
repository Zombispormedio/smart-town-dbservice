package routes

import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "github.com/Zombispormedio/smartdb/controllers"
)



func Set(router *gin.Engine, session *mgo.Session ){


    set:=func(fn func(c *gin.Context, session *mgo.Session), session *mgo.Session) func(c *gin.Context){
        return func(c *gin.Context){       
            fn(c, session)
        } 
    }
    
   
   
    
    Hi:=set(controllers.Hi, session)
    router.GET("/", Hi)
    
    

}
