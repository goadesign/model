package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderBrokerMatchesViewAndDigest(t *testing.T) {
	broker := newRenderBroker()
	results, unregister, err := broker.register("AURA Services", "digest")
	if err != nil {
		t.Fatalf("register render: %v", err)
	}
	defer unregister()

	body := bytes.NewBufferString(
		`{"status":"complete","viewId":"AURA Services","modelDigest":"digest","svg":"<svg/>"}`,
	)
	request := httptest.NewRequest(http.MethodPost, "/headless/result", body)
	response := httptest.NewRecorder()
	broker.handleResult(response, request)
	if response.Code != http.StatusAccepted {
		t.Fatalf("expected accepted result, got %d: %s", response.Code, response.Body.String())
	}
	result := <-results
	if result.ViewID != "AURA Services" || result.ModelDigest != "digest" {
		t.Fatalf("received wrong result: %#v", result)
	}
}

func TestRenderBrokerRejectsUnmatchedResult(t *testing.T) {
	broker := newRenderBroker()
	body := bytes.NewBufferString(
		`{"status":"error","viewId":"Flows","modelDigest":"other","error":"failed"}`,
	)
	request := httptest.NewRequest(http.MethodPost, "/headless/result", body)
	response := httptest.NewRecorder()
	broker.handleResult(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("expected conflict, got %d", response.Code)
	}
}

func TestValidateHeadlessResult(t *testing.T) {
	tests := []struct {
		name   string
		result headlessResult
	}{
		{
			name: "complete without SVG",
			result: headlessResult{
				Status:      "complete",
				ViewID:      "view",
				ModelDigest: "digest",
			},
		},
		{
			name: "error without detail",
			result: headlessResult{
				Status:      "error",
				ViewID:      "view",
				ModelDigest: "digest",
			},
		},
		{
			name: "unknown status",
			result: headlessResult{
				Status:      "running",
				ViewID:      "view",
				ModelDigest: "digest",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validateHeadlessResult(test.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
