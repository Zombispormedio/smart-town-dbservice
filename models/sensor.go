package models

import (
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/lib/struts"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type Sensor struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	NodeID string        `bson:"node_id" json:"node_id"`

	DisplayName   string        `bson:"display_name"  json:"display_name"`
	DeviceName    string        `bson:"device_name"  json:"device_name"`
	Description   string        `bson:"description"  json:"description"`
	SensorGrid    bson.ObjectId `bson:"sensor_grid" json:"sensor_grid"`
	Magnitude     bson.ObjectId `bson:"magnitude" json:"magnitude"`
	Unit          string        `bson:"unit" json:"unit"`
	IsRawData     bool          `bson:"is_raw_data" json:"is_raw_data"`
	RawConversion string        `bson:"raw_conversion" json:"raw_conversion"`

	CreatedBy bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
}

func (sensor *Sensor) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*sensor, reflect.ValueOf(sensor).Elem(), Map, LiteralTag)
}

func SensorCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Sensor")
}

func (sensor *Sensor) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	sensor.FillByMap(obj, "json")
	sensor.ID = bson.NewObjectId()

	newID, _ := uuid.NewV4()
	sensor.NodeID = newID.String()
	sensor.CreatedAt = bson.Now()
	sensor.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorCollection(session)

	InsertError := c.Insert(sensor)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Sensor")

		log.WithFields(log.Fields{
			"message": InsertError.Error(),
		}).Error("SensorInsertError")

		return Error

	}

	sensorGrid := SensorGridCollection(session)

	UpdatingError := sensorGrid.UpdateId(sensor.SensorGrid, bson.M{"$addToSet": bson.M{"sensors": sensor.ID}})

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Pushing  Sensor i SensorGrid: " + sensor.SensorGrid.String())

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      sensor.SensorGrid,
		}).Error("SensorGridSensorInsertError")
	}

	return Error
}

func GetSensors(sensors *[]Sensor, sensorGrid string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)

	iter := c.Find(bson.M{"sensor_grid": bson.ObjectIdHex(sensorGrid)}).Iter()

	IterError := iter.All(sensors)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Sensors")
		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("SensorIteratorError")
	}

	return Error
}

func (sensor *Sensor) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(sensor)
	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding Sensor: " + ID)

		log.WithFields(log.Fields{
			"message": FindingError.Error(),
			"id":      ID,
		}).Warn("SensorByIDError")
	}

	return Error
}

func DeleteSensor(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Sensor: " + ID)
		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Error("SensorRemoveError")
	}

	return Error
}

func (sensor *Sensor) SetTransmissor(ID string, Trans map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"device_name": Trans["device_name"], "description": Trans["description"]}},
		ReturnNew: true,
	}

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensor)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Sensor: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorTransmissorUpdateError")
	}

	return Error
}

func (sensor *Sensor) SetDisplayName(ID string, DisplayName string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)
	change := ChangeOneSet("display_name", DisplayName)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensor)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Sensor: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorDisplayNameUpdateError")
	}

	return Error
}

func (sensor *Sensor) SetMagnitude(ID string, Magnitude map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"magnitude": bson.ObjectIdHex(Magnitude["magnitude"].(string)), "unit": Magnitude["unit"], "is_raw_data": Magnitude["is_raw_data"], "raw_conversion": Magnitude["raw_conversion"]}},
		ReturnNew: true,
	}

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensor)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Sensor: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorMagnitudeUpdateError")
	}

	return Error
}
