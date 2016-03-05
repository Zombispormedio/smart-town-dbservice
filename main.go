package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/Zombispormedio/smartdb/routes"
)

func main() {
	fmt.Println(routes.BuildHello())
	port := os.Getenv("PORT")
    
	if port == "" {
		port="5060"
	}

	router := gin.New()
	router.Use(gin.Logger())
	

	router.GET("/", func(c *gin.Context) {
        
        var msg struct{
            Message string
        }
        
        msg.Message="Hello World"
		c.JSON(200, msg)
	})

	router.Run(":" + port)
}