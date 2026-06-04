package checks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

const defaultHTTPTimeout = 5 * time.Second

// HTTPCheck validates that an HTTP endpoint responds with the expected status code.
type HTTPCheck struct {
	Name           string
	URL            string
	ExpectedStatus int
	Timeout        time.Duration
	Fix            *result.Fix
}

// Run performs an HTTP GET request and checks the status code.
func (c *HTTPCheck) Run(ctx context.Context) result.Result {
	start := time.Now()

	timeout := c.Timeout
	if timeout == 0 {
		timeout = defaultHTTPTimeout
	}

	expectedStatus := c.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return result.Result{
			Name:     c.Name,
			Status:   result.Error,
			Message:  fmt.Sprintf("invalid URL: %s", err),
			Duration: time.Since(start),
			Fix:      c.Fix,
		}
	}

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("%s is unreachable: %s", c.URL, err),
			Hint:     fmt.Sprintf("ensure a service is running at %s", c.URL),
			Duration: duration,
			Fix:      c.Fix,
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != expectedStatus {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("%s returned %d, expected %d", c.URL, resp.StatusCode, expectedStatus),
			Hint:     fmt.Sprintf("expected HTTP %d from %s", expectedStatus, c.URL),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  fmt.Sprintf("%s responded %d", c.URL, resp.StatusCode),
		Duration: duration,
		Fix:      c.Fix,
	}
}
