package osint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// OSINTConfig holds configuration for OSINT modules
type OSINTConfig struct {
	Workers   int
	Timeout   int
	APIKeys   map[string]string
}

var DefaultConfig = OSINTConfig{
	Workers: 10,
	Timeout: 30,
	APIKeys: make(map[string]string),
}

// Result represents the formatted output from an OSINT module
type Result struct {
	Module    string
	Target    string
	Output    string
	Error     error
}

// RunDomain executes domain reconnaissance and returns formatted results
func RunDomain(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runDomainChecks(ctx, target)
	if err != nil {
		return Result{
			Module: "domain",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "domain",
		Target: target,
		Output: formatDomainResult(result),
	}
}

// RunIP executes IP reconnaissance and returns formatted results
func RunIP(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runIPChecks(ctx, target, cfg)
	if err != nil {
		return Result{
			Module: "ip",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "ip",
		Target: target,
		Output: formatIPResult(result),
	}
}

// RunEmail executes email reconnaissance and returns formatted results
func RunEmail(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runEmailChecks(ctx, target)
	if err != nil {
		return Result{
			Module: "email",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "email",
		Target: target,
		Output: formatEmailResult(result),
	}
}

// RunUsername executes username availability checks and returns formatted results
func RunUsername(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runUsernameChecks(ctx, target)
	if err != nil {
		return Result{
			Module: "username",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "username",
		Target: target,
		Output: formatUsernameResult(result),
	}
}

// RunWeb executes web analysis and returns formatted results
func RunWeb(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runWebChecks(ctx, target)
	if err != nil {
		return Result{
			Module: "web",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "web",
		Target: target,
		Output: formatWebResult(result),
	}
}

// RunSSL executes SSL certificate analysis and returns formatted results
func RunSSL(target string, cfg OSINTConfig) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	result, err := runSSLChecks(ctx, target)
	if err != nil {
		return Result{
			Module: "ssl",
			Target: target,
			Error:  err,
		}
	}

	return Result{
		Module: "ssl",
		Target: target,
		Output: formatSSLResult(result),
	}
}

// FormatResult formats an OSINT result for Discord display
func FormatResult(result Result) string {
	if result.Error != nil {
		return fmt.Sprintf("Error: %v", result.Error)
	}
	
	output := fmt.Sprintf("Target: %s\nModule: %s\n\n%s", result.Target, strings.ToUpper(result.Module), result.Output)
	return output
}

// FormatJSON formats an OSINT result as JSON
func FormatJSON(result Result) string {
	data := map[string]interface{}{
		"module":    result.Module,
		"target":    result.Target,
		"timestamp": time.Now().Format(time.RFC3339),
		"output":    result.Output,
	}
	if result.Error != nil {
		data["error"] = result.Error.Error()
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(jsonData)
}
