package checks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestHTTPCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	check := &HTTPCheck{
		Name: "test endpoint",
		URL:  server.URL,
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestHTTPCheck_WrongStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	check := &HTTPCheck{
		Name:           "test endpoint",
		URL:            server.URL,
		ExpectedStatus: http.StatusOK,
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestHTTPCheck_CustomStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	check := &HTTPCheck{
		Name:           "test endpoint",
		URL:            server.URL,
		ExpectedStatus: http.StatusNoContent,
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success for 204, got %s: %s", r.Status, r.Message)
	}
}

func TestHTTPCheck_Unreachable(t *testing.T) {
	check := &HTTPCheck{
		Name:    "unreachable",
		URL:     "http://127.0.0.1:19876/health",
		Timeout: 500 * time.Millisecond,
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestHTTPCheck_InvalidURL(t *testing.T) {
	check := &HTTPCheck{
		Name: "invalid",
		URL:  "://not-a-url",
	}

	r := check.Run(context.Background())
	if r.Status != result.Error {
		t.Errorf("expected error for invalid URL, got %s", r.Status)
	}
}
