package routes

import (
	"github.com/Zombispormedio/smartdb/controllers"
	"github.com/Zombispormedio/smartdb/middleware"
	"github.com/Zombispormedio/smartdb/lib/response"
	"github.com/Zombispormedio/smartdb/consumer"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func Set(router *gin.Engine, session *mgo.Session, consumer *consumer.Consumer) {

	_default := func(fn func(c *gin.Context, session *mgo.Session)) gin.HandlerFunc {
		return func(c *gin.Context) {
			sessionCopy := session.Copy()
			fn(c, sessionCopy)
		}
	}
	
	_defaultAdmin := func(fn func(c *gin.Context, session *mgo.Session)) gin.HandlerFunc {
		return func(c *gin.Context) {
			sessionCopy := session.Copy()
			Error:=middleware.Admin(c, sessionCopy)
			if Error != nil{
				sessionCopy.Close()
				response.Error(c, Error)
				
			}else{
				fn(c, sessionCopy)
			}
		
		}
	}

	router.GET("/", controllers.Hi)
	Status := _default(controllers.Status)
	router.GET("/status", Status)

	api := router.Group("/api")
	{

		oauth := api.Group("/oauth")
		{

			register := _default(controllers.Register)
			oauth.POST("/register", middleware.Secret(), middleware.Body(), register)

			login := _default(controllers.Login)
			oauth.POST("/login", middleware.Body(), login)

			whoiam := _defaultAdmin(controllers.Whoiam)
			oauth.GET("/whoiam",  whoiam)

			logout := _defaultAdmin(controllers.Logout)
			oauth.GET("/logout",  logout)

			DisplayName := _defaultAdmin(controllers.SetOauthDisplayName)
			oauth.PUT("/display_name", middleware.Body(), DisplayName)

			Email := _defaultAdmin(controllers.SetOauthEmail)
			oauth.PUT("/email", middleware.Body(), Email)

			Password := _defaultAdmin(controllers.SetOauthPassword)
			oauth.PUT("/password", middleware.Body(), Password)

			Delete := _defaultAdmin(controllers.DeleteOauth)
			oauth.DELETE("",  Delete)

			Invitate := _defaultAdmin(controllers.Invite)
			oauth.POST("/invite",  middleware.Body(), Invitate)

			CheckInvitation := _default(controllers.CheckInvitation)
			oauth.GET("/invitation/:code", CheckInvitation)
			Invitation := _default(controllers.Invitation)
			oauth.POST("/invitation/:code", middleware.Body(), Invitation)

		}

		magnitude := api.Group("/magnitude")
		{
			Create := _defaultAdmin(controllers.CreateMagnitude)
			magnitude.POST("", middleware.Body(), Create)

			WithID := magnitude.Group("/:id")
			{
				ByID := _defaultAdmin(controllers.GetMagnitudeByID)
				WithID.GET("",  ByID)

				Delete := _defaultAdmin(controllers.DeleteMagnitude)
				WithID.DELETE("",  Delete)

				DisplayName := _defaultAdmin(controllers.SetMagnitudeDisplayName)
				WithID.PUT("/display_name",middleware.Body(), DisplayName)

				Type := _defaultAdmin(controllers.SetMagnitudeType)
				WithID.PUT("/type", middleware.Body(), Type)

				Digital := _defaultAdmin(controllers.SetMagnitudeDigitalUnits)
				WithID.PUT("/digital", middleware.Body(), Digital)

				Analog := WithID.Group("analog")
				{
					Add := _defaultAdmin(controllers.AddMagnitudeAnalogUnit)
					Analog.POST("", middleware.Body(), Add)

					Update := _defaultAdmin(controllers.UpdateMagnitudeAnalogUnit)
					Analog.PUT("", middleware.Body(), Update)

					Delete := _defaultAdmin(controllers.DeleteMagnitudeAnalogUnit)
					Analog.DELETE(":analog_id", Delete)

				}

				Conversion := WithID.Group("conversion")
				{
					Add := _defaultAdmin(controllers.AddMagnitudeConversion)
					Conversion.POST("", middleware.Body(), Add)

					Update := _defaultAdmin(controllers.UpdateMagnitudeConversion)
					Conversion.PUT("", middleware.Body(), Update)

					Delete := _defaultAdmin(controllers.DeleteMagnitudeConversion)
					Conversion.DELETE(":conversion_id", Delete)
				}

			}

		}

		magnitudes := api.Group("/magnitudes")
		{
			All := _defaultAdmin(controllers.GetMagnitudes)
			magnitudes.GET("", All)
			Count := _defaultAdmin(controllers.CountMagnitudes)
			magnitudes.GET("/count", Count)

			Ref := _defaultAdmin(controllers.VerifyRefMagnitude)
			magnitudes.GET("/verify/:ref", Ref)

		}

		zone := api.Group("/zone")
		{
			Create := _defaultAdmin(controllers.CreateZone)
			zone.POST("", middleware.Body(), Create)

			WithID := zone.Group("/:id")
			{

				ByID := _defaultAdmin(controllers.GetZoneByID)
				WithID.GET("", ByID)

				Delete := _defaultAdmin(controllers.DeleteZone)
				WithID.DELETE("", Delete)

				DisplayName := _defaultAdmin(controllers.SetZoneDisplayName)
				WithID.PUT("/display_name", middleware.Body(), DisplayName)

				Description := _defaultAdmin(controllers.SetZoneDescription)
				WithID.PUT("/description", middleware.Body(), Description)

				Keywords := _defaultAdmin(controllers.SetZoneKeywords)
				WithID.PUT("/keywords", middleware.Body(), Keywords)

				Shape := _defaultAdmin(controllers.SetZoneShape)
				WithID.PUT("/shape", middleware.Body(), Shape)
				
				Others := _defaultAdmin(controllers.OtherZones)
				WithID.GET("/others", Others)

			}
		}

		zones := api.Group("/zones")
		{
			All := _defaultAdmin(controllers.GetZones)
			zones.GET("", All)
			Count := _defaultAdmin(controllers.CountZones)
			zones.GET("/count", Count)
			Ref := _defaultAdmin(controllers.VerifyRefZone)
			zones.GET("/verify/:ref", Ref)

		}

		sensorGrids := api.Group("/sensor_grids")
		{
			All := _defaultAdmin(controllers.GetSensorGrids)
			sensorGrids.GET("", All)
			Count := _defaultAdmin(controllers.CountSensorGrids)
			sensorGrids.GET("/count", Count)

			Import := _defaultAdmin(controllers.ImportSensorGrids)
			sensorGrids.POST("/import", middleware.Body(), Import)

			Ref := _defaultAdmin(controllers.VerifyRefSensorGrid)
			sensorGrids.GET("/verify/:ref", Ref)

		}

		sensorGrid := api.Group("/sensor_grid")
		{
			Create := _defaultAdmin(controllers.CreateSensorGrid)
			sensorGrid.POST("", middleware.Body(), Create)

			WithID := sensorGrid.Group("/:id")
			{
				ByID := _defaultAdmin(controllers.GetSensorGridByID)
				WithID.GET("", ByID)

				Delete := _defaultAdmin(controllers.DeleteSensorGrid)
				WithID.DELETE("", Delete)

				Secret := _defaultAdmin(controllers.ChangeSensorGridSecret)
				WithID.GET("/secret", Secret)
				
				MQTT := _defaultAdmin(controllers.ChangeSensorGridMQTT)
				WithID.GET("/mqtt", MQTT)

				CommunicationCenter := _defaultAdmin(controllers.SetSensorGridCommunicationCenter)
				WithID.PUT("/communication_center", middleware.Body(), CommunicationCenter)

				DisplayName := _defaultAdmin(controllers.SetSensorGridDisplayName)
				WithID.PUT("/display_name", middleware.Body(), DisplayName)

				Zone := _defaultAdmin(controllers.SetSensorGridZone)
				WithID.PUT("/zone", middleware.Body(), Zone)

				AllowAccess := _defaultAdmin(controllers.AllowAccessSensorGrid)
				WithID.GET("/access", AllowAccess)

				Location := _defaultAdmin(controllers.SetSensorGridLocation)
				WithID.PUT("/location", middleware.Body(), Location)

				sensors := WithID.Group("/sensors")
				{
					AllSensors := _defaultAdmin(controllers.GetSensors)
					sensors.GET("", AllSensors)

					Count := _defaultAdmin(controllers.CountSensors)
					sensors.GET("/count", Count)

					DelSensor := _defaultAdmin(controllers.DeleteSensorByGrid)
					sensors.DELETE("/:sensor_id", DelSensor)
				}

			}
		}

		sensors := api.Group("/sensors")
		{

			Import := _defaultAdmin(controllers.ImportSensors)
			sensors.POST("/import", middleware.Body(), Import)

			Notifications := _defaultAdmin(controllers.SensorsNotifications)
			sensors.GET("/notifications", Notifications)

		}

		sensor := api.Group("/sensor")
		{
			Create := _defaultAdmin(controllers.CreateSensor)
			sensor.POST("", middleware.Body(), Create)

			WithID := sensor.Group("/:id")
			{
				ByID := _defaultAdmin(controllers.GetSensorByID)
				WithID.GET("", ByID)

				Transmissor := _defaultAdmin(controllers.SetSensorTransmissor)
				WithID.PUT("/transmissor", middleware.Body(), Transmissor)

				DisplayName := _defaultAdmin(controllers.SetSensorDisplayName)
				WithID.PUT("/display_name", middleware.Body(), DisplayName)

				Magnitude := _defaultAdmin(controllers.SetSensorMagnitude)
				WithID.PUT("/magnitude", middleware.Body(), Magnitude)

				Fix := _defaultAdmin(controllers.FixSensor)
				WithID.GET("/fix", Fix)
			}

		}

		task := api.Group("/task")
		{
			Create := _defaultAdmin(controllers.CreateTask)
			task.POST("", middleware.Body(), Create)
			All := _defaultAdmin(controllers.GetTasks)
			task.GET("", All)

			WithID := task.Group("/:id")
			{
				Update := _defaultAdmin(controllers.UpdateTask)
				WithID.PUT("", middleware.Body(), Update)
				Delete := _defaultAdmin(controllers.DeleteTask)
				WithID.DELETE("", Delete)
			}

		}

	}

	push := router.Group("/push", middleware.PushService())
	{

		push.GET("/credentials", controllers.PushCredentialsConfig(consumer))

		SensorRegistry := _default(controllers.PushSensorRegistry)
		push.POST("/sensor_grid", middleware.Body(), SensorRegistry)

		CheckSensorGrid := _default(controllers.CheckSensorGrid)
		push.POST("/sensor_grid/check", middleware.Body(), CheckSensorGrid)

	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404,
			gin.H{"error": gin.H{
				"key":         "not allowed",
				"description": "not allowed",
			}})

	})

}
