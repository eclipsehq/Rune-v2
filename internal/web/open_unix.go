//go:build !windows

package web

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenDashboard() {
	url := "http://localhost:8080"
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}

	if cmd != nil {
		if err := cmd.Start(); err != nil {
			fmt.Printf("[WEB] Failed to auto-launch browser: %v\n", err)
		}
	}
}