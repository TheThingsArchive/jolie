package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

var (
	db *mgo.Session
)

func NewMongoSession() *mgo.Session {
	log.Print("Creating MongoDB connection")
	for {
		log.Printf("Dialling: %s", os.Getenv("MONGODB_URL"))
		s, err := mgo.Dial(fmt.Sprintf("%s:27017", os.Getenv("MONGODB_URL")))
		if err != nil {
			log.Printf("Mongo Connection Error: %s", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			s.SetMode(mgo.Monotonic, true)
			s.SetSocketTimeout(time.Millisecond * 6000)
			s.SetSyncTimeout(time.Millisecond * 6000)
			return s
		}
	}
}
