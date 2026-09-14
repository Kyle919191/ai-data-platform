package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Event struct {
	EventID        string          `json:"event_id"`
	TraceID        string          `json:"trace_id"`
	SessionID      string          `json:"session_id"`
	OrganizationID string          `json:"organization_id"`
	UserID         string          `json:"user_id"`
	EventType      string          `json:"event_type"`
	EventTime      time.Time       `json:"event_time"`
	SchemaVersion  int             `json:"schema_version"`
	Payload        json.RawMessage `json:"payload"`
}

type acceptResponse struct {
	EventID  string `json:"event_id"`
	Accepted bool   `json:"accepted"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if event.EventID == "" {
		event.EventID = fmt.Sprintf("evt_%d", time.Now().UnixNano())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acceptResponse{
		EventID:  event.EventID,
		Accepted: true,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/v1/events", eventsHandler)

	addr := ":8080"
	fmt.Println("ingestion-gateway listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
