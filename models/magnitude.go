package models

import (
	"fmt"
	"time"

	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Conversion struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Operation   string        `bson:"operation" json:"operation"`
	UnitA       bson.ObjectId `bson:"unitA" json:"unitA"`
	UnitB       bson.ObjectId `bson:"unitB" json:"unitB"`
}

type Digital struct {
	ID  bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ON  string        `bson:"on" json:"on"`
	OFF string        `bson:"off" json:"off"`
}

type Analog struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Symbol      string        `bson:"symbol" json:"symbol"`
}

type Magnitude struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Type        string        `bson:"type" json:"type"`
	AnalogUnits []Analog      `bson:"analog_units" json:"analog_units"`
	DigitalUnit Digital       `bson:"digital_units" json:"digital_units"`
	Conversions []Conversion  `bson:"conversions" json:"conversions"`
	CreatedBy   bson.ObjectId `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}



type ListMagnitudeItem struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Type        string        `bson:"type" json:"type"`
	CreatedBy   bson.ObjectId `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}





func MagnitudeCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Magnitude")
}


func (magnitude *Magnitude) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	magnitude.DisplayName = obj["display_name"].(string)
	magnitude.Type = obj["type"].(string)
	magnitude.CreatedAt = bson.Now()
	magnitude.CreatedBy = bson.ObjectIdHex(userID)

	c := MagnitudeCollection(session)

	InsertError := c.Insert(magnitude)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting")
		fmt.Println(InsertError)
	}

	return Error
}

func GetMagnitudes(magnitudes *[]ListMagnitudeItem, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	iter := c.Find(nil).Select(bson.M{"units": 0, "conversions": 0}).Iter()

	IterError := iter.All(magnitudes)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Magnitudes")
		fmt.Println(IterError)
	}

	return Error
}

func DeleteMagnitude(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Magnitude: " + ID)
		fmt.Println(RemoveError)
	}

	return Error
}

func (magnitude *Magnitude) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := MagnitudeCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(magnitude)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding Magnitude: " + ID)
		fmt.Println(FindingError)
	}

	return Error
}

func (magnitude *Magnitude) SetDisplayName(ID string, DisplayName string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)
	change := ChangeOneSet("display_name", DisplayName)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Magnitude " + ID)
		fmt.Println(UpdatingError)
	}

	return Error

}

func (magnitude *Magnitude) SetType(ID string, Type string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)
	change := ChangeOneSet("type", Type)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Magnitude " + ID)
		fmt.Println(UpdatingError)
	}

	return Error

}

func ChangeOneSet(key string, value interface{}) mgo.Change {
	return mgo.Change{
		Update:    bson.M{"$set": bson.M{key: value}},
		ReturnNew: true,
	}
}

func (magnitude *Magnitude) SetDigitalUnits(ID string, units interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	change := ChangeOneSet("digital_units", units)
	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Magnitude Digital Unit " + ID)
		fmt.Println(UpdatingError)
	}

	return Error

}

func (magnitude *Magnitude) AddAnalogUnit(ID string,  unit interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)
    
    preAnalog:=unit.(map[string]interface{})
    
    analog:=Analog{ID: bson.NewObjectId(), DisplayName: preAnalog["display_name"].(string)}
    

	change := mgo.Change{
		Update:    bson.M{"$addToSet": bson.M{"analog_units": analog}},
		ReturnNew: true,
	}
	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Pushing  AnalogUnit" + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}

func (magnitude *Magnitude) UpdateAnalogUnit(ID string, unit interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError
	/*c := MagnitudeCollection(session)
    analogID:=unit.(map[string]interface{})["id"]
	    change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"analog_units.$": unit}},
			ReturnNew: true,
		}

	    _, UpdatingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID), "analog_units._id": bson.ObjectIdHex(analogID.(string))}).Apply(change, &magnitude)

		if UpdatingError != nil {
			Error = utils.BadRequestError("Error Updating  AnalogUnit" + ID)
			fmt.Println(UpdatingError)
		}*/

	return Error

}
