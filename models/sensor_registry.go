package models

import (
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorRegistry struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Sensor     bson.ObjectId `bson:"node_id" json:"node_id"`
	SensorGrid bson.ObjectId `bson:"sensor_grid"    json:"sensor_grid"`
	Value      float64       `bson:"value" json:"value"`
	Date       time.Time     `bson:"date"    json:"date"`
}

func SensorRegistryCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("SensorRegistry")
}

func PushSensorData(Packet []map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	RegistryC := SensorRegistryCollection(session)
	SensorGridC := SensorGridCollection(session)
	SensorC := SensorCollection(session)

	for _, value := range Packet {
		ClientID := value["client_id"].(string)

		Data := utils.SliceInterfaceToSliceMap(value["data"])

		sensorGrid := SensorGrid{}

		SensorGridFindingError := SensorGridC.Find(bson.M{"client_id": ClientID}).One(&sensorGrid)

		if SensorGridFindingError != nil {
			log.WithFields(log.Fields{
				"message":   SensorGridFindingError.Error(),
				"client_id": ClientID,
			}).Warn("SensorGridPushFindingError")
			return utils.BadRequestError("Error Finding SensorGrid in Push: " + ClientID)
		}

		for _, plainSensor := range Data {

			NodeID := plainSensor["node_id"].(string)
			SensorValue := plainSensor["value"].(string)
			SensorDate := plainSensor["date"].(string)

			sensor := Sensor{}

			SensorFindingError := SensorC.Find(bson.M{"node_id": NodeID}).One(&sensor)

			if SensorFindingError != nil {
				log.WithFields(log.Fields{
					"message": SensorFindingError.Error(),
					"node_id": NodeID,
				}).Warn("SensorPushFindingError")
				return utils.BadRequestError("Error Finding Sensor in Push: " + NodeID)
			}

			SensorIsInSensorGrid := utils.ContainsObjectID(sensorGrid.Sensors, sensor.ID)

			if !SensorIsInSensorGrid {
				log.WithFields(log.Fields{
					"message": "Sensor not found in Sensor Grid",
					"id":      sensor.ID,
				}).Warn("SensorPushFindingError")
				return utils.BadRequestError("Error Sensor not found in Sensor Grid: " + sensor.ID.String())
			}

			unixDate, _ := strconv.ParseInt(SensorDate, 10, 64)
			Date := time.Unix(unixDate, 0)

			registry := &SensorRegistry{}

			if utils.Notify(Date) == false {
				registry.Sensor = sensor.ID
				registry.SensorGrid = sensorGrid.ID
				registry.Value, _ = strconv.ParseFloat(SensorValue, 64)
				registry.Date = Date

				InsertError := RegistryC.Insert(registry)

				if InsertError != nil {
					log.WithFields(log.Fields{
						"message": InsertError.Error(),
					}).Warn("SensorPushInsertError")
					return utils.BadRequestError("Error Push Insert")
				}
			}else{
				UpdateError:=SensorC.Update(bson.M{"_id": sensor.ID}, bson.M{"$set":bson.M{"notify":true}})
				
				if UpdateError != nil {
					log.WithFields(log.Fields{
						"message": UpdateError.Error(),
					}).Warn("SensorPushUpdateError")
					return utils.BadRequestError("Error Push Update Notify")
				}
			}

		}

	}

	return Error
}
