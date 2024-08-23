package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080/*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Max-Age", "3600")

}

func handleData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	dataType := r.URL.Query().Get("type")
	data := map[string]string{"message": fmt.Sprintf("Hello from the CORS server! Data type: %s", dataType)}
	json.NewEncoder(w).Encode(data)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	category := r.URL.Query().Get("category")
	id := r.URL.Query().Get("id")
	data := map[string]string{
		"message":  "Info from the CORS server!",
		"category": category,
		"id":       id,
	}
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/api/data", handleData)
	http.HandleFunc("/api/info", handleInfo)

	fmt.Println("CORS Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
