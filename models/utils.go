package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func NextID(c *mgo.Collection) (int, error) {
	var result int
	container := bson.M{}

	Error := c.Find(nil).Select(bson.M{"ref": 1}).Sort("-ref").One(&container)

	if container["ref"] != nil {
		result = container["ref"].(int) + 1
	} else {
		result = 1
	}
	if Error != nil {
		if Error.Error() == "not found" {
			Error = nil
		}

	}

	return result, Error

}

func GetIDbyRef(refStr string, c *mgo.Collection)(bson.ObjectId, error){
		var result bson.ObjectId
	container := bson.M{}
	Ref,_:=strconv.Atoi(refStr)
	Error := c.Find(bson.M{"ref":Ref}).Select(bson.M{"_id": 1}).One(&container)
	
	result=container["_id"].(bson.ObjectId)
	
	return result, Error
}
