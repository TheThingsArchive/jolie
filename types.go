package main

import (
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name,omitempty" bson:"name,omitempty"`
}

type Device struct {
	Id  bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	EUI string        `json:"eui,omitempty" bson:"eui,omitempty"`
}
