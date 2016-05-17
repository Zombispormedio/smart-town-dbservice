package utils

import "github.com/gin-gonic/gin"


func Queries(c *gin.Context, keys []string)map[string]string{
    values:=map[string]string{}
    
    for _,k := range keys{
        
        values[k]=c.Query(k)
        
        
    }
    
    return values
}