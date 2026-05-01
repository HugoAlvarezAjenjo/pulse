package fix

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

const defaultFixTimeout = 120 * time.Second

// Executor handles running fix commands with user confirmation.
type Executor struct {
	Timeout time.Duration
	// Stdin allows overriding the input reader for testing.
	Stdin *os.File
}

// New creates a new fix Executor.
func New() *Executor {
	return &Executor{
		Timeout: defaultFixTimeout,
		Stdin:   os.Stdin,
	}
}

// PromptAndRun asks the user for confirmation, then executes the fix command.
// Returns true if the fix was executed.
func (e *Executor) PromptAndRun(ctx context.Context, r result.Result) (bool, error) {
	if r.Fix == nil || r.Fix.Run == "" {
		return false, nil
	}

	fmt.Printf("  Run fix? (Y/n) ")

	reader := bufio.NewReader(e.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("reading input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		return false, nil
	}

	fmt.Println()
	return true, e.Run(ctx, r.Fix.Run)
}

// Run executes a fix command, streaming output to stdout/stderr.
func (e *Executor) Run(ctx context.Context, command string) error {
	timeout := e.Timeout
	if timeout == 0 {
		timeout = defaultFixTimeout
	}

	fixCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := shellCommandContext(fixCtx, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if fixCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("fix timed out after %s", timeout)
		}
		return fmt.Errorf("fix command failed: %w", err)
	}

	return nil
}

func shellCommandContext(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}
