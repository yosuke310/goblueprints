package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/joeshaw/envdecode"
	"github.com/yosuke310/goblueprints/chapter7/meander"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var env struct {
		APIKey string `env:"SP_GOOGLE_APIKEY,required"`
	}
	if err := envdecode.Decode(&env); err != nil {
		log.Fatalln(err)
	}
	meander.APIKey = env.APIKey

	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}
