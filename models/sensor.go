package models

import (
	"reflect"
	"strconv"
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
	ID            bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	NodeID        string        `bson:"node_id" json:"node_id"`
	Ref           int           `bson:"ref,omitempty" json:"ref"`
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

	LastSync time.Time `bson:"last_sync"    json:"last_sync"`

	Notify bool `bson:"notify"    json:"notify"`
}

func (sensor *Sensor) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*sensor, reflect.ValueOf(sensor).Elem(), Map, LiteralTag)
}

func SensorCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Sensor")
}

func SearchSensorQuery(search string) bson.M {
	or := []bson.M{
		bson.M{"display_name": bson.M{"$regex": search, "$options":"i"}},
		bson.M{"ref": bson.M{"$regex": search, "$options":"i"}},
		bson.M{"description": bson.M{"$regex": search, "$options":"i"}},
		bson.M{"device_name": bson.M{"$regex": search, "$options":"i"}},
		bson.M{"node_id": bson.M{"$regex": search, "$options":"i"}},
		bson.M{"device_name": bson.M{"$regex": search, "$options":"i"}},
	}

	if bson.IsObjectIdHex(search) {
		or = append(or, bson.M{"_id": bson.ObjectIdHex(search)})
		or = append(or, bson.M{"sensor_grid": bson.ObjectIdHex(search)})
		or = append(or, bson.M{"magnitude": bson.ObjectIdHex(search)})
	}

	return bson.M{
		"$or": or,
	}
}

func (sensor *Sensor) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	sensor.FillByMap(obj, "json")
	RefError := sensor.Init(userID, session)

	if RefError != nil {
		log.WithFields(log.Fields{
			"message": RefError.Error(),
		}).Error("SensorRefError")

		return utils.BadRequestError("RefError Sensor: " + RefError.Error())

	}
	c := SensorCollection(session)

	InsertError := c.Insert(sensor)

	if InsertError != nil {
		Error = utils.BadRequestError("InsertError Sensor: " + InsertError.Error())

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

func (sensor *Sensor) Init(userID string, session *mgo.Session) error {
	sensor.ID = bson.NewObjectId()

	newID, _ := uuid.NewV4()
	sensor.NodeID = newID.String()
	sensor.CreatedAt = bson.Now()
	sensor.CreatedBy = bson.ObjectIdHex(userID)
	sensor.LastSync=nil;

	c := SensorCollection(session)

	var RefError error

	sensor.Ref, RefError = NextID(c)

	return RefError

}

func ImportSensors(sensors []map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)
	for _, v := range sensors {
		
		sensor := Sensor{}
		RefError := sensor.Init(userID, session)

		if RefError != nil {
			Error = utils.BadRequestError("RefError Sensor: " + RefError.Error())
			break
		}

		if v["description"] != nil {
			sensor.Description = v["description"].(string)
		}

		if v["display_name"] != nil {
			sensor.DisplayName = v["display_name"].(string)
		}

		if v["device_name"] != nil {
			sensor.DeviceName = v["device_name"].(string)
		}
		
		if v["is_raw_data"] != nil {
			sensor.IsRawData,_=strconv.ParseBool(v["is_raw_data"].(string))
			
			if sensor.IsRawData && v["raw_conversion"]!=nil{
				sensor.RawConversion=v["raw_conversion"].(string)
			}
		}

		if v["magnitude_ref"] != nil && v["magnitude_ref"] != "" {
			magnitudeID, MagnitudeError := GetIDbyRef(v["magnitude_ref"].(string), MagnitudeCollection(session))

			if MagnitudeError != nil {
				Error = utils.BadRequestError("MagnitudeError Sensor: " + MagnitudeError.Error())
				break
			}
			sensor.Magnitude = magnitudeID
		}

		if v["grid_ref"] != nil && v["grid_ref"] != "" {
			gridID, GridError := GetIDbyRef(v["grid_ref"].(string), SensorGridCollection(session))

			if GridError != nil {
				Error = utils.BadRequestError("GridError Sensor: " + GridError.Error())
				break
			}
			sensor.SensorGrid = gridID

			sensorGrid := SensorGridCollection(session)

			UpdatingError := sensorGrid.UpdateId(sensor.SensorGrid, bson.M{"$addToSet": bson.M{"sensors": sensor.ID}})

			if UpdatingError != nil {
				Error = utils.BadRequestError("Error Pushing  Sensor i SensorGrid: " + sensor.SensorGrid.String())

				log.WithFields(log.Fields{
					"message": UpdatingError.Error(),
					"id":      sensor.SensorGrid,
				}).Error("SensorGridSensorInsertError")
			}
		}
		
		InsertError := c.Insert(sensor)

		if InsertError != nil {
			Error = utils.BadRequestError("Error Inserting: " + InsertError.Error() + " Ref: " + string(sensor.Ref))

			break
		}

	}

	return Error
}

func GetSensors(sensors *[]Sensor, sensorGrid string, UrlQuery map[string]string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchSensorGridQuery(search)
	} else {
		query = bson.M{}
	}

	query["sensor_grid"] = bson.ObjectIdHex(sensorGrid)

	var iter *mgo.Iter

	q := c.Find(query).Sort("ref")

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

	IterError := iter.All(sensors)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Sensors")
		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("SensorIteratorError")
	}

	return Error
}

func SensorsNotifications(sensors *[]Sensor, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)
	
	IterError :=  c.Find(bson.M{"notify":true}).Iter().All(sensors)
	
	if IterError != nil {
		Error = utils.BadRequestError("Error Sensors Notifications")
		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("SensorsNotificationsIteratorError")
	}
	
	return Error
}

func FixSensor(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorCollection(session)
	


	UpdatingError:=c.UpdateId(bson.ObjectIdHex(ID), bson.M{"$set": bson.M{"notify":false}})

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Sensor: " + ID)
		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorTransmissorUpdateError")
	}
	
	return Error
}


func CountSensors(sensorGrid string, UrlQuery map[string]string, session *mgo.Session) (int, *utils.RequestError) {
	var Error *utils.RequestError
	var result int
	c := SensorCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchSensorQuery(search)
	} else {
		query = bson.M{}
	}

	query["sensor_grid"] = bson.ObjectIdHex(sensorGrid)

	var CountError error
	result, CountError = c.Find(query).Count()

	if CountError != nil {
		Error = utils.BadRequestError("Error Count Sensors")
		log.WithFields(log.Fields{
			"message": CountError.Error(),
		}).Error("SensorCountError")
	}

	return result, Error
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
