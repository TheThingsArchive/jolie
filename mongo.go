package main

import (
	"fmt"
	"github.com/thethingsnetwork/server-shared"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

const (
	MONGODB_ATTEMPTS = 20
)

type MongoDatabase struct {
	session *mgo.Session
}

func NewMongoDatabase() (*MongoDatabase, error) {
	var err error
	for i := 0; i < MONGODB_ATTEMPTS; i++ {
		uri := os.Getenv("MONGODB_URI")
		var session *mgo.Session
		session, err = mgo.Dial(fmt.Sprintf("%s:27017", uri))
		if err != nil {
			log.Printf("Failed to connect: %s", err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			log.Printf("Connected to %s", err.Error())
			session.SetMode(mgo.Monotonic, true)
			session.SetSocketTimeout(time.Duration(6) * time.Second)
			session.SetSyncTimeout(time.Duration(6) * time.Second)
			return &MongoDatabase{session}, nil
		}
	}
	return nil, err
}

func (db *MongoDatabase) GetApplications() ([]*Application, error) {
	apps := make([]*Application, 0)
	err := db.session.DB("jolie").C("applications").Find(nil).Sort("-_id").Limit(200).All(&apps)
	if err != nil {
		log.Printf("Failed to get applications: %s", err.Error())
		return nil, err
	}
	return apps, nil
}

func (db *MongoDatabase) SaveApplication(app *Application) error {
	err := db.session.DB("jolie").C("applications").Insert(app)
	if err != nil {
		log.Printf("Failed to save application: %s", err.Error())
		return err
	}
	return nil
}

func (db *MongoDatabase) RecordGatewayStatus(status *shared.GatewayStatus) error {
	return nil
}

func (db *MongoDatabase) RecordRxPacket(packet *shared.RxPacket) error {
	return nil
}
