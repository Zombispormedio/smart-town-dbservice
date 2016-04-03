package utils

type RequestError struct{
    Code int
    Message string
}



func BadRequestError(msg string) *RequestError{
    return &RequestError{Code:400, Message:msg}
}


func InterfaceToMapString(obj interface{}) map[string]string{
    return obj.(map[string]string)
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