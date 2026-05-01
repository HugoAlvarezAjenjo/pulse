package checks

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestPortCheck_Open(t *testing.T) {
	// Start a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	check := &PortCheck{
		Name:    "test port",
		Host:    "127.0.0.1",
		Port:    addr.Port,
		Timeout: 2 * time.Second,
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestPortCheck_Closed(t *testing.T) {
	check := &PortCheck{
		Name:    "closed port",
		Host:    "127.0.0.1",
		Port:    19999, // unlikely to be open
		Timeout: 500 * time.Millisecond,
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestPortCheck_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	check := &PortCheck{
		Name:    "timeout port",
		Host:    "192.0.2.1", // non-routable address for timeout testing
		Port:    80,
		Timeout: 50 * time.Millisecond,
	}

	r := check.Run(ctx)
	if r.Status != result.Failure {
		t.Errorf("expected failure on timeout, got %s", r.Status)
	}
}
