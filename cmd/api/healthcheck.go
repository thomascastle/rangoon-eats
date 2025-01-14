package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"env":    app.config.env,
		"status": "available",
	}

	data_JSON, e := json.Marshal(data)
	if e != nil {
		app.logger.Fatal(e, nil)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data_JSON)
}
