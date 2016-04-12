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
	Description string        `bson:"description"  json:"description"`
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

func GetZones(zones *[]Zone, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)

	iter := c.Find(nil).Iter()

	IterError := iter.All(zones)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Zones")
		fmt.Println(IterError)
	}

	return Error
}

func DeleteZone(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Zone: " + ID)
		fmt.Println(RemoveError)
	}

	return Error
}

func (zone *Zone) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := ZoneCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(zone)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding Zone: " + ID)
		fmt.Println(FindingError)
	}

	return Error
}

func (zone *Zone) SetDisplayName(ID string, DisplayName string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)
	change := ChangeOneSet("display_name", DisplayName)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &zone)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Zone: " + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}

func (zone *Zone) SetDescription(ID string, Description string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)
	change := ChangeOneSet("description", Description)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &zone)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Zone: " + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}

func (zone *Zone) SetKeywords(ID string, inKeywords interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError
	
    keywords:=utils.InterfaceToStringArray(inKeywords)
    
    c := ZoneCollection(session)
	change := ChangeOneSet("keywords", keywords)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &zone)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Zone: " + ID)
		fmt.Println(UpdatingError)
	}

	

	return Error
}
