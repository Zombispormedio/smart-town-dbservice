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

type GeoShape struct {
	Type   string      `bson:"type" json:"type"`
	Radius float64     `bson:"radius" json:"radius"`
	Center []float64   `bson:"center" json:"center"`
	Bounds [][]float64 `bson:"bounds" json:"bounds"`
	Paths  [][]float64 `bson:"paths" json:"paths"`
}

func (shape *GeoShape) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*shape, reflect.ValueOf(shape).Elem(), Map, LiteralTag)
}

type Zone struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Keywords    []string      `bson:"keywords" json:"keywords"`
	Center      []float64     `bson:"center" json:"center"`
	Shape       GeoShape      `bson:"shape" json:"shape"`
	CreatedBy   bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
}

func (zone *Zone) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*zone, reflect.ValueOf(zone).Elem(), Map, LiteralTag)
}

func ZoneCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Zone")
}

func (zone *Zone) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	zone.FillByMap(obj, "json")

	zone.CreatedAt = bson.Now()
	zone.CreatedBy = bson.ObjectIdHex(userID)

	c := ZoneCollection(session)

	InsertError := c.Insert(zone)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Zone")
		fmt.Println(InsertError)
	}

	return Error
}
