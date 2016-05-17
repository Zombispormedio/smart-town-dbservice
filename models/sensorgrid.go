package models

import (
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/lib/struts"
	"github.com/Zombispormedio/smartdb/lib/utils"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorGrid struct {
	ID           bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
	Ref          int           `bson:"ref,omitempty" json:"ref"`
	ClientID     string          `bson:"client_id" json:"client_id"`
	ClientSecret string          `bson:"client_secret" json:"client_secret"`
	DisplayName  string          `bson:"display_name"  json:"display_name"`
	DeviceName   string          `bson:"device_name"  json:"device_name"`
	Description  string          `bson:"description"  json:"description"`
	Location     []float64       `bson:"location" json:"location"`
	Zone         bson.ObjectId   `bson:"zone" json:"zone"`
	Sensors      []bson.ObjectId `bson:"sensors" json:"sensors"`
	CreatedBy    bson.ObjectId   `bson:"created_by"    json:"created_by"`
	CreatedAt    time.Time       `bson:"created_at"    json:"created_at"`
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
	newID, _ := uuid.NewV4()
	sensorGrid.ClientID = newID.String()

	sensorGrid.ClientSecret = utils.GenerateSecretToken(47)
	sensorGrid.CreatedAt = bson.Now()
	sensorGrid.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorGridCollection(session)

	var RefError error

	sensorGrid.Ref, RefError = NextID(c)

	if RefError != nil {
		log.WithFields(log.Fields{
			"message": RefError.Error(),
		}).Error("SensorRefError")

		return utils.BadRequestError("RefError SensorGrid: " + RefError.Error())

	}

	InsertError := c.Insert(sensorGrid)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting")

		log.WithFields(log.Fields{
			"message": InsertError.Error(),
		}).Error("SensorGridInsertError")
	}

	return Error
}

func GetSensorGrids(sensorGrids *[]SensorGrid, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	iter := c.Find(nil).Iter()

	IterError := iter.All(sensorGrids)

	if IterError != nil {
		Error = utils.BadRequestError("Error All SensorGrids")

		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("SensorGridIteratorError")
	}

	return Error
}

func (sensorGrid *SensorGrid) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorGridCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(sensorGrid)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding SensorGrid: " + ID)

		log.WithFields(log.Fields{
			"message": FindingError.Error(),
			"id":      ID,
		}).Warn("SensorGridByIDError")
	}

	return Error
}

func DeleteSensorGrid(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	sensorGrid := SensorGrid{}

	FindingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&sensorGrid)

	if FindingError != nil {

		log.WithFields(log.Fields{
			"message": FindingError.Error(),
			"id":      ID,
		}).Warn("SensorGridRemoveFindError")
		return utils.BadRequestError("Error Finding SensorGrid: " + ID)

	}

	count := len(sensorGrid.Sensors)
	var RemovingSensorsError *utils.RequestError
	for index := 0; index < count; index++ {
		sensorID := sensorGrid.Sensors[index].Hex()
		RemovingSensorsError := DeleteSensor(sensorID, session)

		if RemovingSensorsError != nil {
			break
		}

	}

	if RemovingSensorsError != nil {
		return RemovingSensorsError
	}

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing SensorGrid: " + ID)
		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Error("SensorGridRemoveError")
	}

	return Error
}

func (sensorGrid *SensorGrid) ChangeSecret(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorGridCollection(session)
	secret := utils.GenerateSecretToken(47)
	change := ChangeOneSet("client_secret", secret)

	_, UpdatingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID), "client_secret": bson.M{"$ne": nil}}).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("No Allow Access: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridSecretUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) SetCommunicationCenter(ID string, Comm map[string]interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorGridCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"device_name": Comm["device_name"], "description": Comm["description"]}},
		ReturnNew: true,
	}

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating SensorGrid: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridCommunicationCenterUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) SetDisplayName(ID string, DisplayName string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)
	change := ChangeOneSet("display_name", DisplayName)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating SensorGrid: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridDisplayNameUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) SetZone(ID string, ZoneID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)
	change := ChangeOneSet("zone", bson.ObjectIdHex(ZoneID))

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating SensorGrid: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridZoneUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) AllowAccess(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	var change mgo.Change
	c := SensorGridCollection(session)

	temp := SensorGrid{}
	FoundError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&temp)

	if FoundError != nil {

		log.WithFields(log.Fields{
			"message": FoundError.Error(),
			"id":      ID,
		}).Warn("SensorGridFindError")

		return utils.BadRequestError("SensorGrid Not Found: " + ID)
	}

	if temp.ClientSecret != "" {
		change = ChangeOneSet("client_secret", nil)
	} else {
		secret := utils.GenerateSecretToken(47)
		change = ChangeOneSet("client_secret", secret)
	}

	_, UpdatingError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("No Allow Access: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridAccessUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) SetLocation(ID string, location interface{}, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)
	change := ChangeOneSet("location", location)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &sensorGrid)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating SensorGrid: " + ID)

		log.WithFields(log.Fields{
			"message": UpdatingError.Error(),
			"id":      ID,
		}).Warn("SensorGridLocationUpdateError")
	}

	return Error
}

func (sensorGrid *SensorGrid) UnsetSensor(ID string, sensorID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"sensors": bson.ObjectIdHex(sensorID)}},
		ReturnNew: true,
	}

	_, RemoveError := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).Apply(change, &sensorGrid)
	if RemoveError != nil {
		Error = utils.BadRequestError("Error RemovingSensor" + ID)

		log.WithFields(log.Fields{
			"message": RemoveError.Error(),
			"id":      ID,
		}).Error("SensorGridUnsetSensorUpdateError")

		return Error
	}

	InternalError := DeleteSensor(sensorID, session)
	if InternalError != nil {
		Error = utils.BadRequestError("Error RemovingSensor" + ID)
		log.WithFields(log.Fields{
			"id": ID,
		}).Error("SensorRemoveError")

	}

	return Error
}

func (sensorGrid *SensorGrid) CheckCredentials(ClientID string, ClientSecret string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	FindingError := c.Find(bson.M{"client_id": ClientID}).One(sensorGrid)

	if FindingError != nil {

		log.WithFields(log.Fields{
			"client_id": ClientID,
		}).Error("SensorCheckCredencialsError")

		return utils.NoAuthError("Error CheckCredentials- ClientID not found: " + ClientID)
	}

	if ClientSecret != sensorGrid.ClientSecret {
		Error = utils.BadRequestError("Error CheckCredentials- Bad ClientSecret: " + ClientID)
	}

	return Error
}
