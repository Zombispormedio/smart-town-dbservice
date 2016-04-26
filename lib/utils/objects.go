package utils

import(
    "reflect"
 "gopkg.in/mgo.v2/bson"
)

type RequestError struct{
    Code int
    Message string
}


func BadRequestError(msg string) *RequestError{
    return &RequestError{Code:400, Message:msg}
}


func InterfaceToMap(obj interface{}) map[string]interface{}{
    return obj.(map[string]interface{})
}


type TokenLogin struct{
    
    Token string `json:"token"`
}


func Pick(obj map[string]string, pick_name []string)map[string]string{
    
    var result map[string]string
    
    for i := 0; i < len(pick_name); i++ {
        key:=pick_name[i]
        if obj[key]!=""{
            result[key]=obj[key]
        }
        
    }
    
    return result
}

func InterfaceToStringArray(in interface{}) []string{
    s := reflect.ValueOf(in)
	out := make([]string, s.Len())
    
	for i := 0; i < s.Len(); i++ {
		out[i]=s.Index(i).Interface().(string)
	}
    return out
}


func SliceInterfaceToSliceMap(in interface{})[]map[string]interface{}{
    out := []map[string]interface{}{}

	s := reflect.ValueOf(in)

	for i := 0; i < s.Len(); i++ {
		value:=s.Index(i).Interface().(map[string]interface{})
		out=append(out, value)
		
	}
    return out
}


func ContainsObjectID(s []bson.ObjectId, e bson.ObjectId) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}