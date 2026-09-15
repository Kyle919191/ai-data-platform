package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	event := map[string]any{
		"trace_id":   "trace_gen_1",
		"event_type": "model_call_completed",
		"event_time": time.Now().UTC().Format(time.RFC3339),
		"payload": map[string]any{
			"model":      "model-x",
			"latency_ms": 822,
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	resp, err := http.Post(
		"http://localhost:8080/v1/events",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	reply, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("status=%d body=%s\n", resp.StatusCode, reply)
}
