package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Zombispormedio/smartdb/models"
	"gopkg.in/mgo.v2/bson"
)

func isStringType(Kind string, Type string) bool {
	return Kind == "string" && Type == "string"
}

func isObjectIDType(Type string, TypeValue string) bool {
	return Type == "bson.ObjectId" && TypeValue == "string"
}

func isTimeType(Type string, TypeValue string) bool {
	return Type == "time.Time" && TypeValue == "string"
}

func isStruct(Kind string) bool{
    return Kind=="struct"
}

func SetValue(Worker reflect.Value, Field reflect.StructField, Value interface{},LiteralTag string) {
	InnerField := Worker.FieldByName(Field.Name)
	Kind := InnerField.Kind().String()
    Type :=Field.Type
	TypeStr := Type.String()
	TypeValue := reflect.TypeOf(Value).String()

	fmt.Println("Kind: " + Kind + ",  Type: " + TypeStr + ", TypeValue:" + TypeValue)
	fmt.Println(Value)

	switch {
	case isStringType(Kind, TypeStr):
		InnerField.SetString(Value.(string))
	case isObjectIDType(TypeStr, TypeValue):
		ObjectID := bson.ObjectIdHex(Value.(string))
		ObjectIDValue := reflect.ValueOf(ObjectID)
		InnerField.Set(ObjectIDValue)

	case isTimeType(TypeStr, TypeValue):
		Time, _ := time.Parse(time.RFC3339, Value.(string))
		TimeValue := reflect.ValueOf(Time)
		InnerField.Set(TimeValue)

	case isStruct(Kind):
        
        StructValue:=fillValueByMap(Type, Value, LiteralTag)
        InnerField.Set(StructValue)
        
	}

}

func fillValueByMap(Type reflect.Type, Value interface{}, LiteralTag string) reflect.Value{
    
    Worker:=reflect.New(Type)
    
    Len := Type.NumField()

	for i := 0; i < Len; i++ {
		Inner := Type.Field(i)
		KeyMap := Inner.Tag.Get(LiteralTag)
		MapValue := Map[KeyMap]

		SetValue(Worker, Inner, MapValue, LiteralTag)
	}
    
    return Worker
    
}

func FillByMap(Obj interface{}, Worker reflect.Value, Map map[string]interface{}, LiteralTag string) {

	Type := reflect.TypeOf(Obj)

	Len := Type.NumField()

	for i := 0; i < Len; i++ {
		Inner := Type.Field(i)
		KeyMap := Inner.Tag.Get(LiteralTag)
		Value := Map[KeyMap]

		SetValue(Worker, Inner, Value, LiteralTag)
	}

}

func main() {
	Map := map[string]interface{}{
		"_id":          "57029d57c479e0beff645e9c",
		"display_name": "ghjndsd",
		"type":         "0",
		"analog_units": [...]map[string]interface{}{
			map[string]interface{}{
				"_id":          "57029dacc479e0beff645e9e",
				"display_name": "efgrfd",
				"symbol":       "ojng",
			},
			map[string]interface{}{
				"_id":          "57029e15c479e0beff645e9f",
				"display_name": "dfdf",
				"symbol":       "dfdfdf",
			},
		},
		"digital_units": map[string]interface{}{
			"_id": "57029e1ec479e0beff645ea0",
			"on":  "dfdfdwsd",
			"off": "wdske",
		},
		"conversions": [...]map[string]interface{}{
			map[string]interface{}{
				"_id":          "57029e2ac479e0beff645ea1",
				"display_name": "wedssd",
				"operation":    "wdsddsd",
				"unitA":        "57029e33c479e0beff645ea2",
				"unitB":        "57029e3cc479e0beff645ea3",
			},
			map[string]interface{}{
				"_id":          "57029e47c479e0beff645ea4",
				"display_name": "wsddsdf",
				"symbol":       "wdsdff",
				"unitA":        "5702a6d2c479e0beff645ea8",
				"unitB":        "57029d66c479e0beff645e9d",
			},
		},
		"created_by": "56e31e5018e1b981b1017672",
		"created_at": time.Now().Format(time.RFC3339),
	}

	Struct := models.Magnitude{}

	FillByMap(Struct, reflect.ValueOf(&Struct).Elem(), Map, "json")

	fmt.Println(Struct)

}
