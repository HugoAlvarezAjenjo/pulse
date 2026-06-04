package checks

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

const defaultPortTimeout = 2 * time.Second

// PortCheck validates that a TCP port is accepting connections.
type PortCheck struct {
	Name    string
	Host    string
	Port    int
	Timeout time.Duration
	Fix     *result.Fix
}

// Run attempts a TCP connection to the specified host:port.
func (c *PortCheck) Run(ctx context.Context) result.Result {
	start := time.Now()

	timeout := c.Timeout
	if timeout == 0 {
		timeout = defaultPortTimeout
	}

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	// Use a dialer that respects context cancellation
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	duration := time.Since(start)

	if err != nil {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("%s refused connection", address),
			Hint:     fmt.Sprintf("ensure a service is listening on %s", address),
			Duration: duration,
			Fix:      c.Fix,
		}
	}
	_ = conn.Close()

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  fmt.Sprintf("%s accepting connections", address),
		Duration: duration,
		Fix:      c.Fix,
	}
}
