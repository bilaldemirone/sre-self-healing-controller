package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"math"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ERROR: Simulated application failure")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/cpu-load", cpuLoadHandler)
	log.Println("Dummy App started on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func cpuLoadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: CPU load simulation started")

	end := time.Now().Add(30 * time.Second)
	result := 0.0

	for time.Now().Before(end) {
		for i := 0; i < 1000000; i++ {
			result += math.Sqrt(float64(i))
		}
	}

	log.Printf("INFO: CPU load simulation completed. result=%f", result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "cpu load completed",
	})
}