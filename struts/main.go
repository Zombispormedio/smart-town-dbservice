package struts

import (
	"reflect"
	"time"

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
		if RawValue != "" {
			ObjectID := bson.ObjectIdHex(RawValue.(string))
			Value = reflect.ValueOf(ObjectID)
		}

	case IsTimeType(TypeStr, TypeValueStr):
		Time, _ := time.Parse(time.RFC3339, RawValue.(string))
		Value = reflect.ValueOf(Time)

	case IsStructTypeAndMapKindValue(Kind, TypeValue.Kind()):

		Value = fillStructByMap(TypeField, RawValue.(map[string]interface{}), LiteralTag)

	case IsSliceType(Kind):

		Value = fillSliceByMap(TypeField, RawValue, LiteralTag)

	default:
		Value = reflect.ValueOf(RawValue)

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

func fillStructByMap(Type reflect.Type, Map map[string]interface{}, LiteralTag string) reflect.Value {

	Worker := reflect.New(Type).Elem()

	Len := Type.NumField()

	for i := 0; i < Len; i++ {
		Inner := Type.Field(i)
		KeyMap := Inner.Tag.Get(LiteralTag)
		MapValue := Map[KeyMap]
		if MapValue != nil {
			SetValue(Worker, Inner, MapValue, LiteralTag)
		}

	}

	return Worker

}

func FillByMap(Obj interface{}, Worker reflect.Value, Map map[string]interface{}, LiteralTag string) {

	Type := reflect.TypeOf(Obj)

	Len := Type.NumField()

	for i := 0; i < Len; i++ {
		Inner := Type.Field(i)
		KeyMap := Inner.Tag.Get(LiteralTag)
		MapValue := Map[KeyMap]

		if MapValue != nil {

			SetValue(Worker, Inner, MapValue, LiteralTag)
		}

	}

}
