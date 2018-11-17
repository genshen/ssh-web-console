package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func ServeJSON(w http.ResponseWriter, j interface{}) {
	if j == nil {
		http.Error(w, "empty response data", 400)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
	// for request: json.NewDecoder(res.Body).Decode(&body)
}

func Abort(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func GetQueryInt(r *http.Request, key string, defaultValue int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		log.Println("Error: get params cols error:", err)
		return defaultValue
	}
	return value
}

func GetQueryInt32(r *http.Request, key string, defaultValue uint32) uint32 {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		log.Println("Error: get params cols error:", err)
		return defaultValue
	}
	return uint32(value)
}
