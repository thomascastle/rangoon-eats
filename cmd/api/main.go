package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"status": "available",
		}

		data_JSON, e := json.Marshal(data)
		if e != nil {
			log.Fatal(e)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data_JSON)
	})

	log.Println("server started")
	if e := http.ListenAndServe(":4000", nil); e != nil {
		log.Fatal(e)
	}
}
