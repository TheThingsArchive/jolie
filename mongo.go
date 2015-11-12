package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thethingsnetwork/server-shared"
	"gopkg.in/mgo.v2"
)

const (
	MONGODB_ATTEMPTS = 20
)

type MongoDatabase struct {
	session *mgo.Session
}

func ConnectMongoDatabase() (Database, error) {
	var err error
	for i := 0; i < MONGODB_ATTEMPTS; i++ {
		uri := os.Getenv("MONGODB_URI")
		var session *mgo.Session
		session, err = mgo.Dial(fmt.Sprintf("%s:27017", uri))
		if err != nil {
			log.Printf("Failed to connect to %s: %s", uri, err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			log.Printf("Connected to %s", uri)
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

func (db *MongoDatabase) Handle(queues *shared.ConsumerQueues) {
	for {
		select {
		case status := <-queues.GatewayStatuses:
			err := db.session.DB("jolie").C("gateway_statuses").Insert(status)
			log.Printf("Inserted gateway status %#v", status)
			if err != nil {
				//TODO report error and requeue packet
				log.Printf("Failed to save status: %s", err.Error())
			}
		case packet := <-queues.RxPackets:
			err := db.session.DB("jolie").C("rx_packets").Insert(packet)
			log.Printf("Inserted RX packet %#v", packet)
			if err != nil {
				//TODO report error and requeue packet
				log.Printf("Failed to save packet: %s", err.Error())
			}
		}
	}
}

func (db *MongoDatabase) RecordGatewayStatus(status *shared.GatewayStatus) error {
	return nil
}

func (db *MongoDatabase) RecordRxPacket(packet *shared.RxPacket) error {
	return nil
}
