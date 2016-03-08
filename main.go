package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/Zombispormedio/smartdb/config"
    "github.com/Zombispormedio/smartdb/routes"

)



func main() {

    router := gin.New()
    config.ServerConfig(router)

    session:=config.SessionDB()


    if session == nil{
        panic("MongodbSession Fault")
    }else{
        defer session.Close()
    }


    routes.Set(router, session) 

    port := os.Getenv("PORT")

    if port == "" {
        port="5060"
    }

    router.Run(":" + port)
}