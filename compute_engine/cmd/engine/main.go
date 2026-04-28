package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mtepenner/orbital-debris-tracker/compute_engine/internal/conjunction"
	"github.com/mtepenner/orbital-debris-tracker/compute_engine/internal/sgp4"
)

type PredictionRequest struct {
	Objects        []sgp4.CatalogObject `json:"objects"`
	HorizonMinutes int                  `json:"horizon_minutes"`
}

type PredictionResponse struct {
	GeneratedAt  time.Time                  `json:"generated_at"`
	Objects      []sgp4.State               `json:"objects"`
	Conjunctions []conjunction.Conjunction  `json:"conjunctions"`
}

func main() {
	port := getenvInt("COMPUTE_ENGINE_PORT", 7001)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"status": "ok", "service": "compute-engine"})
	})
	mux.HandleFunc("/predict", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payload PredictionRequest
		_ = json.NewDecoder(request.Body).Decode(&payload)
		if len(payload.Objects) == 0 {
			payload.Objects = sgp4.SampleObjects()
		}
		if payload.HorizonMinutes <= 0 {
			payload.HorizonMinutes = 45
		}

		now := time.Now().UTC().Add(time.Duration(payload.HorizonMinutes) * time.Minute)
		states := make([]sgp4.State, 0, len(payload.Objects))
		for _, object := range payload.Objects {
			states = append(states, sgp4.Propagate(object, now))
		}

		response := PredictionResponse{
			GeneratedAt:  now,
			Objects:      states,
			Conjunctions: conjunction.Evaluate(states, 80, now),
		}
		writeJSON(writer, http.StatusOK, response)
	})

	log.Printf("compute engine listening on http://127.0.0.1:%d", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
