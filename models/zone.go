package models

import (
	"reflect"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/lib/struts"
	"github.com/Zombispormedio/smartdb/lib/utils"
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
	Ref         int           `bson:"ref,omitempty" json:"ref"`
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

func SearchZoneQuery(search string) bson.M {
	or := []bson.M{
		bson.M{"display_name": bson.M{"$regex": search}},
		bson.M{"ref": bson.M{"$regex": search}},
		bson.M{"description": bson.M{"$regex": search}},
		bson.M{"keywords": bson.M{"$regex": search}},
		bson.M{"shape.type": bson.M{"$regex": search}},
		bson.M{"center": search},
		bson.M{"shape.radius": search},
		bson.M{"shape.center": search},
		bson.M{"shape.bounds": search},
		bson.M{"shape.paths": search},
	}

	if bson.IsObjectIdHex(search) {
		or = append(or, bson.M{"_id": bson.ObjectIdHex(search)})

	}

	return bson.M{
		"$or": or,
	}
}

func (zone *Zone) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	zone.FillByMap(obj, "json")

	zone.CreatedAt = bson.Now()
	zone.CreatedBy = bson.ObjectIdHex(userID)

	c := ZoneCollection(session)

	var RefError error

	zone.Ref, RefError = NextID(c)

	if RefError != nil {
		log.WithFields(log.Fields{
			"message": RefError.Error(),
		}).Error("ZoneRefError")

		return utils.BadRequestError("RefError Zone: " + RefError.Error())

	}

	InsertError := c.Insert(zone)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Zone")

		log.WithFields(log.Fields{
			"message": InsertError.Error(),
		}).Error("ZoneInsertError")
	}

	return Error
}

func GetZones(zones *[]Zone, UrlQuery map[string]string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchZoneQuery(search)
	}

	var iter *mgo.Iter

	q := c.Find(query)

	if UrlQuery["p"] != "" {
		p, _ := strconv.Atoi(UrlQuery["p"])
		s := 10
		if UrlQuery["s"] != "" {
			s, _ = strconv.Atoi(UrlQuery["s"])
		}

		skip := p * s

		iter = q.Skip(skip).Limit(s).Iter()

	} else {
		iter = q.Iter()
	}

	IterError := iter.All(zones)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Zones")
		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("ZoneIteratorError")
	}

	return Error
}

func CountZones(UrlQuery map[string]string, session *mgo.Session) (int, *utils.RequestError) {
	var Error *utils.RequestError
	var result int
	c := ZoneCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchZoneQuery(search)
	}

	var CountError error
	result, CountError = c.Find(query).Count()

	if CountError != nil {
		Error = utils.BadRequestError("Error Count Zones")
		log.WithFields(log.Fields{
			"message": CountError.Error(),
		}).Error("ZoneCountError")
	}

	return result, Error
}

func DeleteZone(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := ZoneCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Zone: " + ID)

		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Error("ZoneRemoveError")
	}
	sensorGrid := SensorGridCollection(session)

	sensorGrid.Update(bson.M{"zone": bson.ObjectIdHex(ID)}, bson.M{"$unset": bson.M{"zone": true}})

	return Error
}

func (zone *Zone) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := ZoneCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(zone)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding Zone: " + ID)

		log.WithFields(log.Fields{
			"message": FindingError.Error(),
			"id":      ID,
		}).Warn("ZoneByIDError")
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

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("ZoneDisplayNameUpdateError")
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
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("ZoneDescriptionUpdateError")
	}

	return Error
}

func (zone *Zone) SetKeywords(ID string, inKeywords interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError

	keywords := utils.InterfaceToStringArray(inKeywords)

	c := ZoneCollection(session)
	change := ChangeOneSet("keywords", keywords)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &zone)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Zone: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("ZoneKeywordsUpdateError")
	}

	return Error
}

func (zone *Zone) SetShape(ID string, inShape interface{}, inCenter interface{}, session *mgo.Session) *utils.RequestError {

	var Error *utils.RequestError

	shapeMap := utils.InterfaceToMap(inShape)
	shape := GeoShape{}
	shape.FillByMap(shapeMap, "json")

	c := ZoneCollection(session)
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"shape": shape, "center": inCenter}},
		ReturnNew: true,
	}

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &zone)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Zone: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("ZoneShapeUpdateError")
	}

	return Error
}
