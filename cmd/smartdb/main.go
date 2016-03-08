package main

import (
  
	"os"
	"github.com/gin-gonic/gin"
    "github.com/Zombispormedio/smartdb/config"
    
)



func main() {
    
 
	port := os.Getenv("PORT")
    
	if port == "" {
		port="5060"
	}
    
	router := gin.New()
    config.ServerConfig(router)
	

	router.GET("/", func(c *gin.Context) {
        
        var msg struct{
            Message string
        }
        
        msg.Message="Hello World"
		c.JSON(200, msg)
	})

	router.Run(":" + port)
}