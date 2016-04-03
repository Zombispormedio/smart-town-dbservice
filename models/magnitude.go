package models

import (
	"fmt"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Unit struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	DisplayName string        `bson:"display_name"`
	Operation   string        `bson:"operation"`
	UnitA       bson.ObjectId `bson:"unitA"`
	UnitB       bson.ObjectId `bson:"unitB"`
}

type Conversion struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	DisplayName string        `bson:"display_name"`
	Symbol      string        `bson:"symbol"`
	Meaning     []string      `bson:"meaning"`
}

type Magnitude struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	DisplayName string        `bson:"display_name"`
	Type        string        `bson:"type"`
	Units       []Unit        `bson:"units"`
	Conversions []Conversion  `bson:"conversions"`
}

func GetMagnitudeCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Magnitude")
}

func (magnitude *Magnitude) New(obj map[string]string, session *mgo.Session) *utils.RequestError {
	var error *utils.RequestError

	magnitude.DisplayName = obj["display_name"]
	magnitude.Type = obj["type"]

	c := GetMagnitudeCollection(session)

	InsertError := c.Insert(magnitude)

	if InsertError != nil {
		error = utils.BadRequestError("Error Inserting")
		fmt.Println(InsertError)
	}

	return error
}
