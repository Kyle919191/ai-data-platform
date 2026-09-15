package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func validEvent() Event {
	return Event{
		TraceID:   "trace_1",
		EventType: "model_call_completed",
		EventTime: time.Date(2026, 9, 14, 18, 30, 22, 0, time.UTC),
		Payload:   json.RawMessage(`{"model":"model-x"}`),
	}
}

func TestValidateEventOK(t *testing.T) {
	if err := validateEvent(validEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEventMissingTraceID(t *testing.T) {
	event := validEvent()
	event.TraceID = ""

	err := validateEvent(event)
	if err == nil { // this should fail. it doesn't fail(err == nil): smth is wrong
		t.Fatal("expected an error")
	}
	if err.Error() != "trace_id is required" {
		t.Fatalf("got %q", err.Error())
	}
}

func TestEventsHandlerMissingTraceID(t *testing.T) {
	body := `{"event_type":"model_call_completed","event_time":"2026-09-14T18:30:22Z","payload":{"model":"model-x"}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(body)) // build fake request
	rec := httptest.NewRecorder()                                                      // recorder to record in memory

	eventsHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rec.Code)
	}
}
