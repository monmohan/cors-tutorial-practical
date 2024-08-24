package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Custom-Header-1, X-Custom-Header-2")
	(*w).Header().Set("Access-Control-Max-Age", "5")

	if r.Method == "OPTIONS" {
		log.Printf("Preflight request received at %s", time.Now().Format(time.RFC3339))
		log.Printf("Request URL: %s", r.URL.String())
		log.Printf("Request Headers: %v", r.Header)
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func handleData(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Request received for /api/data at %s", time.Now().Format(time.RFC3339))
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Request Headers: %v", r.Header)

	dataType := r.URL.Query().Get("type")
	if dataType == "" {
		dataType = "default"
	}

	data := map[string]string{
		"message":   fmt.Sprintf("Hello from the CORS server! Data type: %s", dataType),
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"header1":   r.Header.Get("X-Custom-Header-1"),
		"header2":   r.Header.Get("X-Custom-Header-2"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/api/data", handleData)

	fmt.Println("CORS Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
