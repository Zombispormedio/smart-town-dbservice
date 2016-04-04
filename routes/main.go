package routes

import (
	"github.com/Zombispormedio/smartdb/controllers"
	"github.com/Zombispormedio/smartdb/middleware"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func Set(router *gin.Engine, session *mgo.Session) {

	_default := func(fn func(c *gin.Context, session *mgo.Session)) gin.HandlerFunc {
		return func(c *gin.Context) {
			sessionCopy := session.Copy()
			fn(c, sessionCopy)
		}
	}

	router.GET("/", controllers.Hi)

	api := router.Group("/api")
	{

		oauth := api.Group("/oauth")
		{

			register := _default(controllers.Register)
			oauth.POST("/register", middleware.Secret(), middleware.Body(), register)

			login := _default(controllers.Login)
			oauth.POST("/login", middleware.Body(), login)

			whoiam := _default(controllers.Whoiam)
			oauth.GET("/whoiam", middleware.Admin(session.Copy()), whoiam)

			logout := _default(controllers.Logout)
			oauth.GET("/logout", middleware.Admin(session.Copy()), logout)
		}

		magnitude := api.Group("/magnitude")
		{
			Create := _default(controllers.CreateMagnitude)
			magnitude.POST("", middleware.Admin(session.Copy()), middleware.Body(), Create)

			All := _default(controllers.GetMagnitudes)
			magnitude.GET("", middleware.Admin(session.Copy()), All)

			WithID := magnitude.Group("/:id")
			{
				ByID := _default(controllers.GetMagnitudeByID)
				WithID.GET("", middleware.Admin(session.Copy()), ByID)

				DisplayName := _default(controllers.SetMagnitudeDisplayName)
				WithID.PUT("/display_name", middleware.Admin(session.Copy()), middleware.Body(), DisplayName)

				Type := _default(controllers.SetMagnitudeType)
				WithID.PUT("/type", middleware.Admin(session.Copy()), middleware.Body(), Type)

				Digital := _default(controllers.SetMagnitudeDigitalUnits)
				WithID.PUT("/digital", middleware.Admin(session.Copy()), middleware.Body(), Digital)

				Analog := WithID.Group("analog")
				{
					Add := _default(controllers.AddMagnitudeAnalogUnit)
					Analog.POST("", middleware.Admin(session.Copy()), middleware.Body(), Add)

					Update := _default(controllers.UpdateMagnitudeAnalogUnit)
					Analog.PUT("/:analog_id", middleware.Admin(session.Copy()), middleware.Body(), Update)

				}

				Delete := _default(controllers.DeleteMagnitude)
				WithID.DELETE("", middleware.Admin(session.Copy()), Delete)
			}

		}

	}

	sensor := router.Group("/sensor")
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
