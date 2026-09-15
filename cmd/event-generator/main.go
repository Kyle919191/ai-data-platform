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

func sendEvent(traceID string) error {
	event := map[string]any{
		"trace_id":   traceID,
		"event_type": "model_call_completed",
		"event_time": time.Now().UTC().Format(time.RFC3339),
		"payload": map[string]any{
			"model":      "model-x",
			"latency_ms": 822,
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"http://localhost:8080/v1/events",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reply, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status=%d body=%s\n", resp.StatusCode, reply)
	return nil
}
func main() {
	for i := 0; i < 10; i++ {
		traceID := fmt.Sprintf("trace_gen_%d", i)
		if err := sendEvent(traceID); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
