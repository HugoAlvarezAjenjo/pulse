package checks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// FileCheck validates that a file exists at the specified path.
type FileCheck struct {
	Name string
	Path string
	Fix  *result.Fix
}

// Run checks if the file exists.
func (c *FileCheck) Run(_ context.Context) result.Result {
	start := time.Now()

	info, err := os.Stat(c.Path)
	duration := time.Since(start)

	if err != nil {
		if os.IsNotExist(err) {
			return result.Result{
				Name:     c.Name,
				Status:   result.Failure,
				Message:  fmt.Sprintf("file not found: %s", c.Path),
				Hint:     fmt.Sprintf("create the file at %s", c.Path),
				Duration: duration,
				Fix:      c.Fix,
			}
		}
		return result.Result{
			Name:     c.Name,
			Status:   result.Error,
			Message:  fmt.Sprintf("error checking file: %s", err),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	if info.IsDir() {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("path is a directory, not a file: %s", c.Path),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  fmt.Sprintf("found %s", c.Path),
		Duration: duration,
		Fix:      c.Fix,
	}
}
