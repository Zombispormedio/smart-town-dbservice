package models

import (
	"reflect"
	"time"
    "fmt"
	"github.com/Zombispormedio/smartdb/utils"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/struts"
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`

	DisplayName string `bson:"display_name"  json:"display_name"`
	Webhook     string `bson:"webhook"  json:"webhook"`
	Frequency   string `bson:"frequency"  json:"frequency"`

	CreatedBy bson.ObjectId `bson:"created_by"    json:"created_by"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
}

func (task *Task) FillByMap(Map map[string]interface{}, LiteralTag string) {
	struts.FillByMap(*task, reflect.ValueOf(task).Elem(), Map, LiteralTag)
}

func TaskCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("Task")
}


func (task *Task) New(obj map[string]interface{}, userID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	task.FillByMap(obj, "json")
   
	task.CreatedAt = bson.Now()
	task.CreatedBy = bson.ObjectIdHex(userID)

	c := TaskCollection(session)

	InsertError := c.Insert(task)

	if InsertError != nil {
		Error = utils.BadRequestError("Error Inserting Task")
		fmt.Println(InsertError)
	}

	return Error
}



func GetTasks(tasks *[]Task, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := TaskCollection(session)

	iter := c.Find(nil).Iter()

	IterError := iter.All(tasks)

	if IterError != nil {
		Error = utils.BadRequestError("Error All Tasks")
		fmt.Println(IterError)
	}

	return Error
}




func (task *Task) Update( ID string, obj map[string]interface{},session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError

	c := TaskCollection(session)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"display_name": obj["display_name"], "webhook":obj["webhook"], "frequency":obj["frequency"]}},
		ReturnNew: true,
	}

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &task)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating Task: " + ID)
		fmt.Println(UpdatingError)
	}


	return Error
}


func DeleteTask(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := TaskCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing Task: " + ID)
		fmt.Println(RemoveError)
	}

	return Error
}