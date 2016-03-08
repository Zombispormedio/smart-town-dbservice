package config

import (
   
	"github.com/gin-gonic/gin"

)

func ServerConfig(router *gin.Engine){
    router.Use(gin.Logger())

}


