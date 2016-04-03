package models

import (
	"fmt"
	"time"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Unit struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Operation   string        `bson:"operation" json:"operation"`
	UnitA       bson.ObjectId `bson:"unitA"`
	UnitB       bson.ObjectId `bson:"unitB"`
}

type Conversion struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Symbol      string        `bson:"symbol" json:"symbol"`
	Meaning     []string      `bson:"meaning" json:"meaning"`
}

type Magnitude struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Type        string        `bson:"type" json:"type"`
	Units       []Unit        `bson:"units" json:"units"`
	Conversions []Conversion  `bson:"conversions" json:"conversions"`
	CreatedBy   bson.ObjectId `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}

type ListMagnitudeItem struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Type        string        `bson:"type" json:"type"`
	CreatedBy   bson.ObjectId `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}

func MagnitudeCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Magnitude")
}

func (magnitude *Magnitude) New(obj map[string]string, userID string, session *mgo.Session) *utils.RequestError {
	var error *utils.RequestError

	magnitude.DisplayName = obj["display_name"]
	magnitude.Type = obj["type"]
	magnitude.CreatedAt = bson.Now()
	magnitude.CreatedBy = bson.ObjectIdHex(userID)

	c := MagnitudeCollection(session)

	InsertError := c.Insert(magnitude)

	if InsertError != nil {
		error = utils.BadRequestError("Error Inserting")
		fmt.Println(InsertError)
	}

	return error
}

func AllMagnitudes(magnitudes *[]ListMagnitudeItem, session *mgo.Session) *utils.RequestError {
	var ReqError *utils.RequestError
	c := MagnitudeCollection(session)

	iter := c.Find(nil).Select(bson.M{"units": 0, "conversions": 0}).Iter()

	IterError := iter.All(magnitudes)

	if IterError != nil {
		ReqError = utils.BadRequestError("Error All Magnitudes")
		fmt.Println(IterError)
	}

	return ReqError
}

func DelMagnitude(id string, session *mgo.Session) *utils.RequestError {
	var ReqError *utils.RequestError
	c := MagnitudeCollection(session)
    
    RemoveError:=c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
    
    if RemoveError != nil {
		ReqError = utils.BadRequestError("Error Removing Magnitude: "+id)
		fmt.Println(RemoveError)
	}

	return ReqError
}
