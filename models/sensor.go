package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/struts"
	"github.com/Zombispormedio/smartdb/utils"
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
	newID, _ := uuid.NewV4()
	sensor.NodeID = newID.String()

	
	sensor.CreatedAt = bson.Now()
	sensor.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorCollection(session)

	InsertError := c.Insert(sensor)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Sensor")
		fmt.Println(InsertError)
	}

	return Error
}
