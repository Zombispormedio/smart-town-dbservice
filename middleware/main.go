package middleware

import(

    "os"
    "strings"
    "github.com/gin-gonic/gin"
    //  "gopkg.in/mgo.v2"

    "github.com/Zombispormedio/smartdb/response"


)

func Sensor() gin.HandlerFunc{
    return func(c *gin.Context){ 
        c.Next()
    }
}

func Secret() gin.HandlerFunc{
    return func(c *gin.Context){

        equals := strings.Compare(c.Request.Header.Get("Authorization"), os.Getenv("SMARTDBSECRET"))
        if equals !=0{
            response.Error(c, 403, "No Authorization");
            return;
        }
        c.Next();

    }
}