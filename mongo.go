package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

const (
	MONGODB_ATTEMPTS = 20
)

var (
	db *mgo.Session
)

func NewMongoSession() (*mgo.Session, error) {
	var err error
	for i := 0; i < MONGODB_ATTEMPTS; i++ {
		uri := os.Getenv("MONGODB_URI")
		var s *mgo.Session
		s, err = mgo.Dial(fmt.Sprintf("%s:27017", uri))
		if err != nil {
			log.Printf("Failed to connect: %s", err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			log.Printf("Connected to %s", err.Error())
			s.SetMode(mgo.Monotonic, true)
			s.SetSocketTimeout(time.Millisecond * 6000)
			s.SetSyncTimeout(time.Millisecond * 6000)
			return s, nil
		}
	}
	return nil, err
}
