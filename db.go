package main

import (
	"fmt"
	"log"

	"labix.org/v2/mgo"
)

var getDb func() (*mgo.Database, *mgo.Session)

func initDb(dbHost, dbName string) {
	url := fmt.Sprintf("mongodb://%s/%s", dbHost, dbName)
	mgoSession, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Dialed:", url)

	getDb = func() (*mgo.Database, *mgo.Session) {
		s := mgoSession.Clone()
		return s.DB(dbName), s
	}
}
