package middleware

import (
	"os"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
		
	"github.com/Zombispormedio/smartdb/models"
	"github.com/Zombispormedio/smartdb/lib/response"
)


func Secret() gin.HandlerFunc {
	return func(c *gin.Context) {

		equals := strings.Compare(c.Request.Header.Get("Authorization"), os.Getenv("SMARTDBSECRET"))
		if equals != 0 {
			response.ErrorByString(c, 403, "No Authorization: SecretError")
			return
		}
		c.Next()

	}
}

func Body() gin.HandlerFunc {
	return func(c *gin.Context) {

		var body interface{}

		BindingJSONError := c.BindJSON(&body)

		if BindingJSONError != nil {
			response.ErrorByString(c, 400, "No body in HttpRequest")
			return
		}

		c.Set("body", body)

		c.Next()
	}
}

func Admin(session *mgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")

		if token == "" {
			response.ErrorByString(c, 403, "No Authorization: Empty Token")
			return
		}

		user, err := models.SessionToken(token, session)

		if err != nil {
			response.Error(c, err)
			return
		}

		c.Set("user", user)

		c.Next()

	}
	
}



func PushService() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		fmt.Println(token)

		c.Next()

	}
}



