package response

import(
    "github.com/gin-gonic/gin"

)



type MessageT struct{
    Status int  `json:"status"`
    Message string `json:"message"`
}

type ErrorT struct{
    Status int  `json:"status"`
    Error string `json:"error"`
}

type DataT struct{
    Status int  `json:"status"`
    Data interface{} `json:"data"`
}



func SuccessMessage(c *gin.Context, message string){
    var msg MessageT

    msg.Message=message
    c.JSON(200, msg)
}

func Error(c * gin.Context, code int, error string){
    var msg ErrorT
    
    msg.Status=1
    msg.Error=error
    c.JSON(code, msg)
    c.AbortWithStatus(code)
}

func Success(c *gin.Context, data interface{}){
    var msg DataT

    msg.Data=data
    c.JSON(200, data)
}