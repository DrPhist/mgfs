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
	checkError(err, true)

	log.Println("Dialed:", url)

	getDb = func() (*mgo.Database, *mgo.Session) {
		s := mgoSession.Clone()
		return s.DB(dbName), s
	}
}
