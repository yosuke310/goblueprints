package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/stretchr/graceful"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "Endpoint address")
		mongo = flag.String("mongo", "localhost", "MongoDB address")
	)
	flag.Parse()
	log.Println("MongoDB Connecting...")
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("MongoDB Connect failed:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Web server Starting...")
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("Web server Stopping...")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "Invalid API Key")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withData(d *mgo.Session, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("polls"))
		fn(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
