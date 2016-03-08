package config
import(
    
    "os"
    "log"
    "time"
    "gopkg.in/mgo.v2"
)

func SessionDB() *mgo.Session {

    mongoDBDialInfo:=&mgo.DialInfo{
        Addrs: []string{os.Getenv("SMARTDBHOST")},
        Timeout:60*time.Second,
        Database: os.Getenv("SMARTDB"),
        Username: os.Getenv("SMARTDBUSER"),
        Password: os.Getenv("SMARTDBPASS"),
    }

    session, err:= mgo.DialWithInfo(mongoDBDialInfo)
    if err != nil{
        log.Fatalf("CreateSession: %s\n", err)
        return nil
    }

    session.SetMode(mgo.Monotonic, true)

    
    return session

}


func GetDB(session *mgo.Session) *mgo.Database{
    return session.DB(os.Getenv("SMARTDB"))
}
