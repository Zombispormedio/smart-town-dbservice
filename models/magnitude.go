package models

import (
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/lib/struts"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Conversion struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Operation   string        `bson:"operation" json:"operation"`
	UnitA       bson.ObjectId `bson:"unitA" json:"unitA"`
	UnitB       bson.ObjectId `bson:"unitB" json:"unitB"`
}

func (conversion *Conversion) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*conversion, reflect.ValueOf(conversion).Elem(), Map, LiteralTag)
}

type Digital struct {
	ID  bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ON  string        `bson:"on"    json:"on"`
	OFF string        `bson:"off"   json:"off"`
}

func (digital *Digital) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*digital, reflect.ValueOf(digital).Elem(), Map, LiteralTag)
}

type Analog struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Symbol      string        `bson:"symbol"    json:"symbol"`
}

func (analog *Analog) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*analog, reflect.ValueOf(analog).Elem(), Map, LiteralTag)
}

type Magnitude struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Type        string        `bson:"type"  json:"type"`
	AnalogUnits []Analog      `bson:"analog_units"  json:"analog_units"`
	DigitalUnit Digital       `bson:"digital_units" json:"digital_units"`
	Conversions []Conversion  `bson:"conversions"   json:"conversions"`
	CreatedBy   bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
}

func (magnitude *Magnitude) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*magnitude, reflect.ValueOf(magnitude).Elem(), Map, LiteralTag)
}

func MagnitudeCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Magnitude")
}

func (magnitude *Magnitude) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	magnitude.FillByMap(obj, "json")

	magnitude.CreatedAt = bson.Now()
	magnitude.CreatedBy = bson.ObjectIdHex(userID)

	c := MagnitudeCollection(session)

	InsertError := c.Insert(magnitude)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Magnitude")
		log.WithFields(log.Fields{
			"message": InsertError.Error(),
		}).Error("MagnitudeInsertError")
	}

	return Error
}

func GetMagnitudes(magnitudes *[]Magnitude, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	iter := c.Find(nil).Select(bson.M{"units": 0, "conversions": 0}).Iter()

	IterError := iter.All(magnitudes)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Magnitudes")
		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("MagnitudeInteratorError")
	}

	return Error
}

func DeleteMagnitude(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Magnitude: " + ID)

		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Error("MagnitudeRemoveError")
	}

	sensor := SensorCollection(session)

	sensor.Update(bson.M{"magnitude": bson.ObjectIdHex(ID)}, bson.M{"$unset": bson.M{"magnitude": true, "unit": true}})

	return Error
}

func (magnitude *Magnitude) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := MagnitudeCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(magnitude)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding Magnitude: " + ID)
		log.WithFields(log.Fields{
			"message": FindingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeFindByIDError")
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
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeUpdateDisplayNameError")
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
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeUpdateTypeError")
	}

	return Error

}

func ChangeOneSet(key string, value interface{}) mgo.Change {
	return mgo.Change{
		Update:    bson.M{"$set": bson.M{key: value}},
		ReturnNew: true,
	}
}

func (magnitude *Magnitude) SetDigitalUnits(ID string, units map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	digital := Digital{}
	digital.FillByMap(units, "json")

	if !digital.ID.Valid() {
		digital.ID = bson.NewObjectId()
	}
	change := ChangeOneSet("digital_units", digital)
	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Magnitude Digital Unit " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeUpdateDigitalUnitsError")
	}

	return Error

}

func (magnitude *Magnitude) AddAnalogUnit(ID string, unit map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	analog := Analog{}
	analog.FillByMap(unit, "json")

	if !analog.ID.Valid() {
		analog.ID = bson.NewObjectId()
	}

	change := mgo.Change{
		Update:    bson.M{"$addToSet": bson.M{"analog_units": analog}},
		ReturnNew: true,
	}
	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Pushing  AnalogUnit" + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeInsertAnalogUnitError")
	}

	return Error
}

func (magnitude *Magnitude) UpdateAnalogUnit(ID string, unit map[string]interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	analog := Analog{}
	analog.FillByMap(unit, "json")

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"analog_units.$": analog}},
		ReturnNew: true,
	}

	_, UpdatingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID), "analog_units._id": analog.ID}).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating  AnalogUnit" + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeUpdateAnalogUnitError")
	}

	return Error

}

func (magnitude *Magnitude) DeleteAnalogUnit(ID string, analogID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"analog_units": bson.M{"_id": bson.ObjectIdHex(analogID)}}},
		ReturnNew: true,
	}

	_, RemoveError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).Apply(change, &magnitude)
	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing AnalogUnit" + ID)
		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Warn("MagnitudeRemoveAnalogUnitError")
	}

	return Error

}

func (magnitude *Magnitude) AddConversion(ID string, unit map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	conversion := Conversion{}
	conversion.FillByMap(unit, "json")

	if !conversion.ID.Valid() {
		conversion.ID = bson.NewObjectId()
	}

	change := mgo.Change{
		Update:    bson.M{"$addToSet": bson.M{"conversions": conversion}},
		ReturnNew: true,
	}
	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Pushing  Conversion: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeInsertConversionError")
	}

	return Error
}

func (magnitude *Magnitude) UpdateConversion(ID string, inConversion map[string]interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	conversion := Conversion{}
	conversion.FillByMap(inConversion, "json")

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"conversions.$": conversion}},
		ReturnNew: true,
	}

	_, UpdatingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID), "conversions._id": conversion.ID}).Apply(change, &magnitude)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating  Conversion" + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("MagnitudeUpdateConversionError")
	}

	return Error

}

func (magnitude *Magnitude) DeleteConversion(ID string, conversionID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := MagnitudeCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"conversions": bson.M{"_id": bson.ObjectIdHex(conversionID)}}},
		ReturnNew: true,
	}

	_, RemoveError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).Apply(change, &magnitude)
	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Conversion" + ID)
		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Warn("MagnitudeRemoveConversionError")
	}

	return Error

}
