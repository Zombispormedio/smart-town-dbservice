package models

import(

    "errors"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/Zombispormedio/smartdb/config"

    "fmt"

)

type OAuth struct{
    Email string `bson:"email"`
    Password string `bson:"password"`
    Token []string  `bson: "token"`
}


func GetOAuthCollection(session *mgo.Session) *mgo.Collection{
    return config.GetDB(session).C("OAuth")
}

func (oauth *OAuth) New(obj  map[string]string, session *mgo.Session) error{


    if(obj["email"]=="" || obj["password"]==""){
        return errors.New("Empty Params")
    }

    c:=GetOAuthCollection(session)

    result := OAuth{}

    err := c.Find(bson.M{"email": obj["email"]}).One(&result)

    fmt.Println(err)
    fmt.Println(result)



    return nil

}