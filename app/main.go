package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Debug endpoint — menampilkan APP_DEBUG_KEY dari environment
	mux.HandleFunc("/v1/debug", func(w http.ResponseWriter, r *http.Request) {
		debugKey := os.Getenv("APP_DEBUG_KEY")
		if debugKey == "" {
			debugKey = "NOT_SET"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"app":           "zenvault",
			"app_debug_key": debugKey,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 ZenVault API running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
