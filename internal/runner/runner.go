package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

func RunShell(ctx context.Context, command string) (string, error) {
	var executable string
	var args []string
	if runtime.GOOS == "windows" {
		executable = "cmd"
		args = []string{"/C", command}
	} else {
		executable = "sh"
		args = []string{"-c", command}
	}

	cmd := exec.CommandContext(ctx, executable, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return output.String(), fmt.Errorf("command failed: %w", err)
	}
	return output.String(), nil
}
