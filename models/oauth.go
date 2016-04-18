package models

import (
	"os"
	"time"
	"fmt"
	"github.com/Zombispormedio/smartdb/config"
	"github.com/Zombispormedio/smartdb/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/nu7hatch/gouuid"
)

type OAuth struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	DisplayName string        `bson:"display_name"  json:"display_name"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	Token       []string      `bson:"token"`
	Invitation  string        `bson:"invitation"  json:"invitation"`
}

type Profile struct {
	Email       string `bson:"email" json:"email"`
	DisplayName string `bson:"display_name"  json:"display_name"`
}

func GetOAuthCollection(session *mgo.Session) *mgo.Collection {
	return config.GetDB(session).C("OAuth")
}

func getOAuthByEmail(email string, session *mgo.Session) (OAuth, error) {

	c := GetOAuthCollection(session)

	result := OAuth{}

	error := c.Find(bson.M{"email": email}).One(&result)

	return result, error
}

func insertOAuth(oauth *OAuth, session *mgo.Session) error {
	c := GetOAuthCollection(session)

	InsertError := c.Insert(oauth)

	return InsertError
}

func (oauth *OAuth) Register(obj map[string]interface{}, session *mgo.Session) *utils.RequestError {

	if obj["email"] == "" || obj["password"] == "" {
		return utils.BadRequestError("Empty Params")
	}

	_, FoundError := getOAuthByEmail(obj["email"].(string), session)

	if FoundError == nil {
		return utils.BadRequestError("User Exists")
	}

	oauth.Email = obj["email"].(string)
	password := []byte(obj["password"].(string))

	hashedPassword, EncryptError := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if EncryptError != nil {
		return utils.BadRequestError("Error With Password")
	}
	oauth.Password = string(hashedPassword)

	InsertError := insertOAuth(oauth, session)

	if InsertError != nil {
		return utils.BadRequestError("Error Inserting")
	}

	return nil

}

func Login(obj map[string]interface{}, session *mgo.Session) (*utils.TokenLogin, *utils.RequestError) {

	if obj["email"] == "" || obj["password"] == "" {
		return nil, utils.BadRequestError("Empty Params")
	}

	oauth, FoundError := getOAuthByEmail(obj["email"].(string), session)

	if FoundError != nil {
		return nil, utils.BadRequestError("User Not Exists")
	}

	tempHashPass := []byte(oauth.Password)
	tempStrPass := []byte(obj["password"].(string))

	PasswordError := bcrypt.CompareHashAndPassword(tempHashPass, tempStrPass)

	if PasswordError != nil {
		return nil, utils.BadRequestError("Password Not Valid")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = oauth.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	secret := []byte(os.Getenv("SMARTDBSECRET"))
	tokenString, TokenGeneratedError := token.SignedString(secret)

	if TokenGeneratedError != nil {
		return nil, utils.BadRequestError("TokenError")

	}

	collection := GetOAuthCollection(session)

	tokenQuery := bson.M{"token": tokenString}
	UpdateError := collection.Update(bson.M{"_id": oauth.ID}, bson.M{"$addToSet": tokenQuery})

	if UpdateError != nil {
		return nil, utils.BadRequestError("Updating Token Array Error")

	}

	return &utils.TokenLogin{Token: tokenString}, nil
}

func SessionToken(tokenString string, session *mgo.Session) (interface{}, *utils.RequestError) {

	token, ParsingTokenError := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		return []byte(os.Getenv("SMARTDBSECRET")), nil
	})

	if ParsingTokenError != nil {
		return nil, utils.BadRequestError("Parsing Token Error")
	}

	oauthID := token.Claims["id"]

	collection := GetOAuthCollection(session)

	result := OAuth{}

	FindingError := collection.Find(bson.M{"_id": bson.ObjectIdHex(oauthID.(string)), "token": tokenString}).One(&result)

	if FindingError != nil {
		return nil, utils.BadRequestError("Session Not Exists")
	}

	return oauthID, nil

}

func Logout(token string, userID string, session *mgo.Session) *utils.RequestError {

	collection := GetOAuthCollection(session)

	updateQuery := bson.M{"token": token}
	UpdateError := collection.UpdateId(bson.ObjectIdHex(userID), bson.M{"$pull": updateQuery})

	if UpdateError != nil {
		return utils.BadRequestError("Updating Problem")
	}

	return nil
}

func GetProfile(userID string, profile *Profile, session *mgo.Session) *utils.RequestError {

	collection := GetOAuthCollection(session)

	FindingError := collection.Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(profile)

	if FindingError != nil {
		return utils.BadRequestError("User Not Exists")
	}

	

	return nil
}



func (profile *Profile) SetDisplayName(ID string, DisplayName string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := GetOAuthCollection(session)
	change := ChangeOneSet("display_name", DisplayName)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &profile)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating profile: " + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}

func (profile *Profile) SetEmail(ID string, Email string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := GetOAuthCollection(session)
	change := ChangeOneSet("email", Email)

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &profile)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating profile: " + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}


func (profile *Profile) SetPassword(ID string, Password string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := GetOAuthCollection(session)
	
	password := []byte(Password)

	hashedPassword, EncryptError := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	
	if EncryptError != nil {
		return utils.BadRequestError("Error With Password")
	}
	
	change := ChangeOneSet("password", string(hashedPassword))

	_, UpdatingError := c.FindId(bson.ObjectIdHex(ID)).Apply(change, &profile)

	if UpdatingError != nil {
		Error = utils.BadRequestError("Error Updating profile: " + ID)
		fmt.Println(UpdatingError)
	}

	return Error
}



func DeleteOauth(ID string, session *mgo.Session) *utils.RequestError {
	var Error *utils.RequestError
	c := GetOAuthCollection(session)

	RemoveError := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})

	if RemoveError != nil {
		Error = utils.BadRequestError("Error Removing User: " + ID)
		fmt.Println(RemoveError)
	}

	return Error
}

func Invite (Email string, session *mgo.Session) (string,*utils.RequestError) {
	var Error *utils.RequestError
	code:=""
	
	c := GetOAuthCollection(session)
	
	count, _:=c.Find(bson.M{"email":Email}).Count()
	
	if count!=0{
		return code, utils.BadRequestError("Error Email Exists")
	}
	
	newID, _ := uuid.NewV4()
	code = newID.String()
	guest:=OAuth{}
	guest.Email=Email
	guest.Invitation=code;
	
	InsertError := insertOAuth(&guest, session)

	if InsertError != nil {
		Error=utils.BadRequestError("Error Inserting Invitation")
	}

	
	return code, Error
}

func CheckInvitation (code string, session *mgo.Session) (bool,*utils.RequestError) {
	var checked bool
	var Error *utils.RequestError
	c := GetOAuthCollection(session)
	
	count, _:=c.Find(bson.M{"invitation":code}).Count()
	
	if count==0{
		return checked, utils.BadRequestError("Error Invitation not Exists")
	}
	
	checked=true
	
	return checked, Error
}

func AcceptInvitation (Code string, Password string, session *mgo.Session) *utils.RequestError{
	var Error *utils.RequestError
	
	c := GetOAuthCollection(session)
	
	oauth:=OAuth{}
	
	FindingError:=c.Find(bson.M{"invitation":Code}).One(&oauth)
	fmt
	if FindingError != nil {
		fmt.Println(FindingError)
		return utils.BadRequestError("Error Finding Guest")
	}
	
	password := []byte(Password)

	hashedPassword, EncryptError := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	
	if EncryptError != nil {
		fmt.Println(EncryptError)
		return utils.BadRequestError("Error With Password")
	}
	
	
	UpdateError:=c.UpdateId(oauth.ID, bson.M{"$set":bson.M{"password": string(hashedPassword)}, "$unset": bson.M{"invitation":true}})
	
	if UpdateError != nil {
		fmt.Println(UpdateError)
		Error=utils.BadRequestError("Error Updating Invitation")
	}
	
	return Error
}