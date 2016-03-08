package config

/*
 
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
    }
    
    defer session.Close()
    
    session.SetMode(mgo.Monotonic, true)
    
    
    c:= session.DB(mongoDBDialInfo.Database).C("people")
    err = c.Insert(&Person{"Ale", "+5584155451"},
                   &Person{"Cla", "+5258452545"})
    
    if err != nil{
        log.Fatal(err)
    }
    result:=Person{}
    err= c.Find(bson.M{"name": "Ale"}).One(&result)
    
    if err != nil{
        log.Fatal(err)
    }
    
    fmt.Println("Phone: ", result.Phone)*/