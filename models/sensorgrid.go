package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/struts"
	"github.com/Zombispormedio/smartdb/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorGrid struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Zone        bson.ObjectId `bson:"zone" json:"zone"`
	CreatedBy   bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
}

func (sensorGrid *SensorGrid) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*sensorGrid, reflect.ValueOf(sensorGrid).Elem(), Map, LiteralTag)
}

func SensorGridCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("SensorGrid")
}

func (sensorGrid *SensorGrid) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	sensorGrid.FillByMap(obj, "json")

	sensorGrid.CreatedAt = bson.Now()
	sensorGrid.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorGridCollection(session)

	InsertError := c.Insert(sensorGrid)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting")
		fmt.Println(InsertError)
	}

	return Error
}
