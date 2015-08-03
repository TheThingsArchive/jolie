package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"ApplicationsIndex",
		"GET",
		"/applications",
		ApplicationsIndexHandler,
	},
	Route{
		"ApplicationsCreate",
		"POST",
		"/applications",
		ApplicationsCreateHandler,
	},
	Route{
		"DevicesIndex",
		"GET",
		"/devices",
		DevicesIndexHandler,
	},
}

func JSONResponse(status int, data interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to The Jolie API")
}

func ApplicationsIndexHandler(w http.ResponseWriter, r *http.Request) {
	res := make([]Application, 0)
	err := db.DB("jolie").C("applications").Find(nil).Sort("-_id").Limit(200).All(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(http.StatusOK, res, w)
}

func ApplicationsCreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	app := new(Application)
	err := decoder.Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Insert Datas
	err = db.DB("jolie").C("applications").Insert(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(http.StatusCreated, app, w)
}

func DevicesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Listing Devices registered with the things network")
}

func Api() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
