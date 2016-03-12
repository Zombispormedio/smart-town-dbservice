package models

import(
    "os"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "golang.org/x/crypto/bcrypt"
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/Zombispormedio/smartdb/config"
    "github.com/Zombispormedio/smartdb/utils"

    "fmt"

)

type OAuth struct{
    ID   bson.ObjectId `bson:"_id,omitempty"`
    Email string `bson:"email"`
    Password string `bson:"password"`
    Token []string  `bson: "token"`
}


func GetOAuthCollection(session *mgo.Session) *mgo.Collection{
    return config.GetDB(session).C("OAuth")
}

func getOAuthByEmail(email string, session *mgo.Session) (OAuth, error){

    c:=GetOAuthCollection(session)

    result := OAuth{}

    error := c.Find(bson.M{"email": email}).One(&result)

 
    return result, error
}

func insertOAuth(oauth *OAuth, session *mgo.Session) error{
    c:=GetOAuthCollection(session)

    InsertError := c.Insert(oauth)

    return InsertError
}

func (oauth *OAuth) Register(obj  map[string]string, session *mgo.Session)  *utils.RequestError{


    if(obj["email"]=="" || obj["password"]==""){
        return utils.BadRequestError("Empty Params")
    }

    _, FoundError :=getOAuthByEmail(obj["email"], session)

    if FoundError == nil{
        return utils.BadRequestError("User Exists")
    }
    
  
    oauth.Email=obj["email"];
    password:=[]byte(obj["password"])

    hashed_password, EncryptError:= bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

    if EncryptError != nil{
        return utils.BadRequestError("Error With Password")
    }
    oauth.Password=string(hashed_password)

    InsertError := insertOAuth(oauth, session)

    if InsertError != nil{
        return utils.BadRequestError("Error Inserting")
    }

    return nil

}

func Login( obj  map[string]string, session *mgo.Session) (*utils.TokenLogin, *utils.RequestError){


    if(obj["email"]=="" || obj["password"]==""){
        return nil, utils.BadRequestError("Empty Params")
    }

    oauth, FoundError :=getOAuthByEmail(obj["email"], session)

    if FoundError != nil{
        return nil, utils.BadRequestError("User Not Exists")
    }
    
    tempHashPass:=[]byte(oauth.Password)
    tempStrPass:=[]byte(obj["password"])

    PasswordError := bcrypt.CompareHashAndPassword(tempHashPass, tempStrPass)
    fmt.Println(PasswordError)
    if  PasswordError != nil{
        return nil, utils.BadRequestError("Password Not Valid")
    }


    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims["id"] =oauth.ID
    token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    
    secret:=[]byte(os.Getenv("SMARTDBSECRET"))
    tokenString, TokenGeneratedError := token.SignedString(secret)
    fmt.Println(TokenGeneratedError)

    if  TokenGeneratedError != nil{
        return nil, utils.BadRequestError("TokenError")
    }


    return &utils.TokenLogin{Token:tokenString}, nil
}

