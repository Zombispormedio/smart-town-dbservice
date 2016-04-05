package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Zombispormedio/smartdb/models"
	"gopkg.in/mgo.v2/bson"
)

func IsStringType(Kind string, Type string) bool {
	return Kind == "string" && Type == "string"
}

func IsObjectIDType(Type string, TypeValue string) bool {
	return Type == "bson.ObjectId" && TypeValue == "string"
}

func IsTimeType(Type string, TypeValue string) bool {
	return Type == "time.Time" && TypeValue == "string"
}

func IsStructTypeAndMapKindValue(Kind string, KindValue reflect.Kind) bool {
	return Kind == "struct" && KindValue == reflect.Map
}

func IsSliceType(Kind string) bool {
	return Kind == "slice"
}

func MakeValue(KindField reflect.Kind, TypeField reflect.Type, RawValue interface{}, LiteralTag string) reflect.Value {
	var Value reflect.Value

	Kind := KindField.String()
	TypeStr := TypeField.String()

	TypeValue := reflect.TypeOf(RawValue)
	TypeValueStr := TypeValue.String()

	switch {
	case IsStringType(Kind, TypeStr):
		Value = reflect.ValueOf(RawValue)
	case IsObjectIDType(TypeStr, TypeValueStr):
		ObjectID := bson.ObjectIdHex(RawValue.(string))
		Value = reflect.ValueOf(ObjectID)

	case IsTimeType(TypeStr, TypeValueStr):
		Time, _ := time.Parse(time.RFC3339, RawValue.(string))
		Value = reflect.ValueOf(Time)

	case IsStructTypeAndMapKindValue(Kind, TypeValue.Kind()):

		Value = fillValueByMap(TypeField, RawValue.(map[string]interface{}), LiteralTag)

	case IsSliceType(Kind):

		Value = fillSliceByMap(TypeField, RawValue, LiteralTag)

	}

	return Value
}

func SetValue(Worker reflect.Value, Field reflect.StructField, Value interface{}, LiteralTag string) {
	InnerField := Worker.FieldByName(Field.Name)
	NewValue := MakeValue(InnerField.Kind(), InnerField.Type(), Value, LiteralTag)

	if NewValue.IsValid() {
		InnerField.Set(NewValue)
	}

}

func fillSliceByMap(Type reflect.Type, RawValue interface{}, LiteralTag string) reflect.Value {

	TempSlice := reflect.ValueOf(RawValue)

	l := TempSlice.Len()

	SliceValue := reflect.MakeSlice(Type, l, l)

	for i := 0; i < l; i++ {
		Elem := TempSlice.Index(i)
		raw := Elem.Interface()
		v := SliceValue.Index(i)

		ElemValue := MakeValue(v.Kind(), v.Type(), raw, LiteralTag)

		v.Set(ElemValue)

	}

	return SliceValue
}

func fillValueByMap(Type reflect.Type, Map map[string]interface{}, LiteralTag string) reflect.Value {

	Worker := reflect.New(Type).Elem()

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
				"operation":    "wdsdff",
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
