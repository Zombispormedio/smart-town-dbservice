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
			
			DisplayName := _default(controllers.SetOauthDisplayName)
			oauth.PUT("/display_name", middleware.Admin(session.Copy()), middleware.Body(), DisplayName)
			
			Email := _default(controllers.SetOauthEmail)
			oauth.PUT("/email", middleware.Admin(session.Copy()), middleware.Body(), Email)
			
			Password := _default(controllers.SetOauthPassword)
			oauth.PUT("/password", middleware.Admin(session.Copy()), middleware.Body(), Password)
			
			Delete := _default(controllers.DeleteOauth)
			oauth.DELETE("", middleware.Admin(session.Copy()), Delete)
			
			Invitate := _default(controllers.Invite)
			oauth.POST("/invite", middleware.Admin(session.Copy()), middleware.Body(), Invitate)
			
			
			CheckInvitation := _default(controllers.CheckInvitation)
			oauth.GET("/invitation/:code",  CheckInvitation)
			Invitation := _default(controllers.Invitation)
			oauth.POST("/invitation/:code", middleware.Body(), Invitation)
			
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

				Delete := _default(controllers.DeleteMagnitude)
				WithID.DELETE("", middleware.Admin(session.Copy()), Delete)

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
					Analog.PUT("", middleware.Admin(session.Copy()), middleware.Body(), Update)

					Delete := _default(controllers.DeleteMagnitudeAnalogUnit)
					Analog.DELETE(":analog_id", middleware.Admin(session.Copy()), Delete)

				}

				Conversion := WithID.Group("conversion")
				{
					Add := _default(controllers.AddMagnitudeConversion)
					Conversion.POST("", middleware.Admin(session.Copy()), middleware.Body(), Add)

					Update := _default(controllers.UpdateMagnitudeConversion)
					Conversion.PUT("", middleware.Admin(session.Copy()), middleware.Body(), Update)

					Delete := _default(controllers.DeleteMagnitudeConversion)
					Conversion.DELETE(":conversion_id", middleware.Admin(session.Copy()), Delete)
				}

			}

		}

		zone := api.Group("/zone")
		{
			Create := _default(controllers.CreateZone)
			zone.POST("", middleware.Admin(session.Copy()), middleware.Body(), Create)

			All := _default(controllers.GetZones)
			zone.GET("", middleware.Admin(session.Copy()), All)
			WithID := zone.Group("/:id")
			{

				ByID := _default(controllers.GetZoneByID)
				WithID.GET("", middleware.Admin(session.Copy()), ByID)

				Delete := _default(controllers.DeleteZone)
				WithID.DELETE("", middleware.Admin(session.Copy()), Delete)

				DisplayName := _default(controllers.SetZoneDisplayName)
				WithID.PUT("/display_name", middleware.Admin(session.Copy()), middleware.Body(), DisplayName)

				Description := _default(controllers.SetZoneDescription)
				WithID.PUT("/description", middleware.Admin(session.Copy()), middleware.Body(), Description)

				Keywords := _default(controllers.SetZoneKeywords)
				WithID.PUT("/keywords", middleware.Admin(session.Copy()), middleware.Body(), Keywords)

				Shape := _default(controllers.SetZoneShape)
				WithID.PUT("/shape", middleware.Admin(session.Copy()), middleware.Body(), Shape)

			}
		}

		sensorGrid := api.Group("/sensor_grid")
		{
			Create := _default(controllers.CreateSensorGrid)
			sensorGrid.POST("", middleware.Admin(session.Copy()), middleware.Body(), Create)
			All := _default(controllers.GetSensorGrids)
			sensorGrid.GET("", middleware.Admin(session.Copy()), All)

			WithID := sensorGrid.Group("/:id")
			{
				ByID := _default(controllers.GetSensorGridByID)
				WithID.GET("", middleware.Admin(session.Copy()), ByID)

				Delete := _default(controllers.DeleteSensorGrid)
				WithID.DELETE("", middleware.Admin(session.Copy()), Delete)

				Secret := _default(controllers.ChangeSensorGridSecret)
				WithID.GET("/secret", middleware.Admin(session.Copy()), Secret)

				CommunicationCenter := _default(controllers.SetSensorGridCommunicationCenter)
				WithID.PUT("/communication_center", middleware.Admin(session.Copy()), middleware.Body(), CommunicationCenter)

				DisplayName := _default(controllers.SetSensorGridDisplayName)
				WithID.PUT("/display_name", middleware.Admin(session.Copy()), middleware.Body(), DisplayName)

				Zone := _default(controllers.SetSensorGridZone)
				WithID.PUT("/zone", middleware.Admin(session.Copy()), middleware.Body(), Zone)

				AllowAccess := _default(controllers.AllowAccessSensorGrid)
				WithID.GET("/access", middleware.Admin(session.Copy()), AllowAccess)

				Location := _default(controllers.SetSensorGridLocation)
				WithID.PUT("/location", middleware.Admin(session.Copy()), middleware.Body(), Location)

				sensors := WithID.Group("/sensors")
				{
					AllSensors := _default(controllers.GetSensors)
					sensors.GET("", middleware.Admin(session.Copy()), AllSensors)

					DelSensor := _default(controllers.DeleteSensorByGrid)
					sensors.DELETE("/:sensor_id", middleware.Admin(session.Copy()), DelSensor)
				}

			}
		}
		sensor := api.Group("/sensor")
		{
			Create := _default(controllers.CreateSensor)
			sensor.POST("", middleware.Admin(session.Copy()), middleware.Body(), Create)

			WithID := sensor.Group("/:id")
			{
				ByID := _default(controllers.GetSensorByID)
				WithID.GET("", middleware.Admin(session.Copy()), ByID)

				Transmissor := _default(controllers.SetSensorTransmissor)
				WithID.PUT("/transmissor", middleware.Admin(session.Copy()), middleware.Body(), Transmissor)

				DisplayName := _default(controllers.SetSensorDisplayName)
				WithID.PUT("/display_name", middleware.Admin(session.Copy()), middleware.Body(), DisplayName)
				
				Magnitude := _default(controllers.SetSensorMagnitude)
				WithID.PUT("/magnitude", middleware.Admin(session.Copy()), middleware.Body(), Magnitude)
			}

		}
		
		task:=api.Group("/task")
		{
			Create := _default(controllers.CreateTask)
			task.POST("", middleware.Admin(session.Copy()), middleware.Body(), Create)
			All := _default(controllers.GetTasks)
			task.GET("", middleware.Admin(session.Copy()), All)
			
			WithID := task.Group("/:id")
			{
				Update := _default(controllers.UpdateTask)
				WithID.PUT("", middleware.Admin(session.Copy()), middleware.Body(), Update)
				Delete := _default(controllers.DeleteTask)
				WithID.DELETE("", middleware.Admin(session.Copy()), Delete)
			}
			
		}

	}
	
	
	push := router.Group("/push", middleware.PushService())
	{
		config:=push.Group("/config")
		{
			config.GET("/credentials", controllers.PushCredentialsConfig)
			
			SensorGrid:=_default(controllers.CheckSensorGrid)
			config.POST("/sensor_grid", middleware.Body(), SensorGrid)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404,
			gin.H{"error": gin.H{
				"key":         "not allowed",
				"description": "not allowed",
			}})

	})

}
