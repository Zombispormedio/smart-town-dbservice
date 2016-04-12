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
    "github.com/nu7hatch/gouuid"
)

type SensorGrid struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ClientID     string        `bson:"client_id,omitempty" json:"client_id"`
	ClientSecret string        `bson:"client_secret",omitempty" json:"client_secret"`
	DisplayName  string        `bson:"display_name"  json:"display_name"`
	Zone         bson.ObjectId `bson:"zone" json:"zone"`
	CreatedBy    bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt    time.Time     `bson:"created_at"    json:"created_at"`
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
    newID, _:=uuid.NewV4()
    sensorGrid.ClientID=newID.String()
    
    sensorGrid.ClientSecret=utils.GenerateSecretToken(47)
	sensorGrid.CreatedAt = bson.Now()
	sensorGrid.CreatedBy = bson.ObjectIdHex(userID)

	c := SensorGridCollection(session)

	InsertError := c.Insert(sensorGrid)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting")
		fmt.Println(InsertError)
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
		fmt.Println(IterError)
	}

	return Error
}

func (sensorGrid *SensorGrid) ByID(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := SensorGridCollection(session)

	FindingError := c.FindId(bson.ObjectIdHex(ID)).One(sensorGrid)

	if FindingError != nil {
		Error = utils.BadRequestError("Error Finding SensorGrid: " + ID)
		fmt.Println(FindingError)
	}

	return Error
}

func DeleteSensorGrid(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := SensorGridCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing SensorGrid: " + ID)
		fmt.Println(RemoveError)
	}

	return Error
}
