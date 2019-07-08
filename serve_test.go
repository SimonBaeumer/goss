package goss

import (
	"github.com/SimonBaeumer/goss/outputs"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHealthHandler_Serve(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	cmdResource := &resource.Command{
		Command:    "echo hello",
		Title:      "echo hello",
		ExitStatus: 0,
	}

	h := HealthHandler{
		Cache:         cache.New(time.Duration(50), time.Duration(50)),
		Outputer:      outputs.GetOutputer("documentation"),
		MaxConcurrent: 1,
		ListenAddr:    "9999",
		ContentType:   "application/json",
		GossMu:        &sync.Mutex{},
		GossConfig: GossConfig{
			Commands: resource.CommandMap{"echo hello": cmdResource},
		},
	}
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check failed!")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Title: echo hello")
	assert.Contains(t, rr.Body.String(), "Command: echo hello: exit-status: matches expectation: [0]")
	assert.Contains(t, rr.Body.String(), "Count: 1, Failed: 0, Skipped: 0")
}
