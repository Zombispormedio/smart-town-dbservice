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

type SensorGrid struct {
	ID           bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
	Ref          int             `bson:"ref,omitempty" json:"ref"`
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

func SearchSensorGridQuery(search string) bson.M {
	or := []bson.M{
		bson.M{"display_name": bson.M{"$regex": search}},
		bson.M{"ref": bson.M{"$regex": search}},
		bson.M{"description": bson.M{"$regex": search}},
		bson.M{"device_name": bson.M{"$regex": search}},
		bson.M{"client_id": bson.M{"$regex": search}},
		bson.M{"location": search},
	}

	if bson.IsObjectIdHex(search) {
		or = append(or, bson.M{"_id": bson.ObjectIdHex(search)})
		or = append(or, bson.M{"zone": bson.ObjectIdHex(search)})
		or = append(or, bson.M{"sensors": bson.ObjectIdHex(search)})
	}

	return bson.M{
		"$or": or,
	}
}

func (sensorGrid *SensorGrid) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	sensorGrid.FillByMap(obj, "json")

	RefError := sensorGrid.Init(userID, session)

	if RefError != nil {
		log.WithFields(log.Fields{
			"message": RefError.Error(),
		}).Error("SensorGridRefError")

		return utils.BadRequestError("RefError SensorGrid: " + RefError.Error())

	}

	c := SensorGridCollection(session)

	InsertError := c.Insert(sensorGrid)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting")

		log.WithFields(log.Fields{
			"message": InsertError.Error(),
		}).Error("SensorGridInsertError")
	}

	return Error
}

func (sensorGrid *SensorGrid) Init(userID string, session *mgo.Session) error {
	newID, _ := uuid.NewV4()
	sensorGrid.ClientID = newID.String()

	sensorGrid.ClientSecret = utils.GenerateSecretToken(47)
	sensorGrid.CreatedAt = bson.Now()
	sensorGrid.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorGridCollection(session)

	var RefError error

	sensorGrid.Ref, RefError = NextID(c)

	return RefError

}

func ImportSensorGrids(grids []map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)
	for _, v := range grids {

		grid := SensorGrid{}
		RefError := grid.Init(userID, session)

		if RefError != nil {
			Error = utils.BadRequestError("RefError SensorGrid: " + RefError.Error())
			break
		}

		if v["description"] != nil {
			grid.Description = v["description"].(string)
		}

		if v["display_name"] != nil {
			grid.DisplayName = v["display_name"].(string)
		}

		if v["device_name"] != nil {
			grid.DeviceName = v["device_name"].(string)
		}

		if v["location_longitude"] != nil && v["location_latitude"] != nil {
			long, _ := strconv.ParseFloat(v["location_longitude"].(string), 64)
			lat, _ := strconv.ParseFloat(v["location_latitude"].(string), 64)
			grid.Location = []float64{
				long,lat, 
			}
		}

		if v["zone_ref"] != nil && v["zone_ref"] !=""{
			ZoneID, ZoneError := GetIDbyRef(v["zone_ref"].(string), ZoneCollection(session))

			if ZoneError != nil {
				Error = utils.BadRequestError("ZoneError SensorGrid: " + ZoneError.Error())
				break
			}
			grid.Zone = ZoneID
		}

		InsertError := c.Insert(grid)

		if InsertError != nil {
			Error = utils.BadRequestError("Error Inserting: " + InsertError.Error() + " Ref: " + string(grid.Ref))

			break
		}

	}

	return Error
}

func GetSensorGrids(sensorGrids *[]SensorGrid, UrlQuery map[string]string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchSensorGridQuery(search)
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

	IterError := iter.All(sensorGrids)

	if IterError != nil {
		Error = utils.BadRequestError("Error All SensorGrids")

		log.WithFields(log.Fields{
			"message": IterError.Error(),
		}).Error("SensorGridIteratorError")
	}

	return Error
}

func CountSensorGrids(UrlQuery map[string]string, session *mgo.Session) (int, *utils.RequestError) {
	var Error *utils.RequestError
	var result int
	c := SensorGridCollection(session)

	var query bson.M

	if UrlQuery["search"] != "" {
		search := UrlQuery["search"]
		query = SearchSensorGridQuery(search)
	}

	var CountError error
	result, CountError = c.Find(query).Count()

	if CountError != nil {
		Error = utils.BadRequestError("Error Count SensorGrids")
		log.WithFields(log.Fields{
			"message": CountError.Error(),
		}).Error("SensorGridCountError")
	}

	return result, Error
}




func VerifyRefSensorGrid(RefStr string, session *mgo.Session) (bool, *utils.RequestError) {
	var Error *utils.RequestError
	var result bool
	c := SensorGridCollection(session)
	
	
	Ref, _:=strconv.Atoi(RefStr)
	
	count,  CountError:=c.Find(bson.M{"ref": Ref}).Count()
	
	if CountError != nil {
		Error = utils.BadRequestError("Error Ref SensorGrid: "+CountError.Error())
		log.WithFields(log.Fields{
			"message": CountError.Error(),
		}).Error("SensorGridRefError")
	}

	if count==1{
		result=true
	}

	return result, Error
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
