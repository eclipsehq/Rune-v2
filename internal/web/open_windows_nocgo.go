//go:build windows && !cgo

package web

import (
	"fmt"
	"os/exec"
)

func OpenDashboard() {
	url := "http://localhost:8080"
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	if err := cmd.Start(); err != nil {
		fmt.Printf("[WEB] Failed to auto-launch browser: %v\n", err)
	}
}