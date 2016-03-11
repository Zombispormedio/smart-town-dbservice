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