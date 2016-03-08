package routes

import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "github.com/Zombispormedio/smartdb/controllers"
    "github.com/Zombispormedio/smartdb/middleware"
)



func Set(router *gin.Engine, session *mgo.Session ){


    _default:=func(fn func(c *gin.Context, session *mgo.Session)) gin.HandlerFunc{
        return func(c *gin.Context){  
            sessionCopy := session.Copy()
            fn(c, sessionCopy)
        } 
    }





    router.GET("/", controllers.Hi)


    api := router.Group("/api")
    {
        api_with_secret:=api.Use(middleware.Secret())

        register:=_default(controllers.Register)
        api_with_secret.POST("/register", register)

    }

    sensor:=router.Group("/sensor")
    sensor.Use(middleware.Sensor())
    {

    }

    router.NoRoute(func(c *gin.Context) {
        c.JSON(404,
               gin.H{"error": gin.H{
                   "key":         "not allowed",
                   "description": "not allowed",
               }})

    })





}
