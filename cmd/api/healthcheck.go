package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": "development",
		"version":     version,
	}

	js, err := json.Marshal(data)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}
