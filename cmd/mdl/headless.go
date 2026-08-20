// This file matches browser render results to the exact CLI request that
// started them. View IDs and model digests prevent one run from accepting
// another run's output.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type (
	headlessResult struct {
		Status      string `json:"status"`
		ViewID      string `json:"viewId"`
		ModelDigest string `json:"modelDigest"`
		SVG         string `json:"svg,omitempty"`
		Error       string `json:"error,omitempty"`
	}

	renderKey struct {
		viewID      string
		modelDigest string
	}

	renderBroker struct {
		mu      sync.Mutex
		waiters map[renderKey]chan headlessResult
	}
)

// newRenderBroker creates an empty result matcher for one mdl svg process.
func newRenderBroker() *renderBroker {
	return &renderBroker{waiters: make(map[renderKey]chan headlessResult)}
}

// register creates the one result channel allowed for a view and model digest.
func (b *renderBroker) register(viewID, modelDigest string) (<-chan headlessResult, func(), error) {
	key := renderKey{viewID: viewID, modelDigest: modelDigest}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.waiters[key]; exists {
		return nil, nil, fmt.Errorf("render %s is already registered", viewID)
	}
	result := make(chan headlessResult, 1)
	b.waiters[key] = result
	unregister := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.waiters, key)
	}
	return result, unregister, nil
}

// handleResult accepts one complete or failed browser result for a registered request.
func (b *renderBroker) handleResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<20))
	decoder.DisallowUnknownFields()
	var result headlessResult
	if err := decoder.Decode(&result); err != nil {
		http.Error(w, fmt.Sprintf("decode render result: %v", err), http.StatusBadRequest)
		return
	}
	if err := validateHeadlessResult(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := renderKey{viewID: result.ViewID, modelDigest: result.ModelDigest}
	b.mu.Lock()
	waiter, exists := b.waiters[key]
	if exists {
		select {
		case waiter <- result:
		default:
			exists = false
		}
	}
	b.mu.Unlock()
	if !exists {
		http.Error(w, "no matching render request", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// validateHeadlessResult requires exactly the fields allowed by each status.
func validateHeadlessResult(result headlessResult) error {
	if result.ViewID == "" || result.ModelDigest == "" {
		return fmt.Errorf("viewId and modelDigest are required")
	}
	switch result.Status {
	case "complete":
		if result.SVG == "" || result.Error != "" {
			return fmt.Errorf("complete result requires svg and forbids error")
		}
	case "error":
		if result.Error == "" || result.SVG != "" {
			return fmt.Errorf("error result requires error and forbids svg")
		}
	default:
		return fmt.Errorf("invalid render status %q", result.Status)
	}
	return nil
}
