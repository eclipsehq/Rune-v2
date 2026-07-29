package osint

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WebResult struct {
	URL           string                 `json:"url"`
	StatusCode    int                    `json:"status_code"`
	StatusText    string                 `json:"status_text"`
	Headers       map[string]string     `json:"headers"`
	Redirects     []string               `json:"redirects"`
	TechStack     []string               `json:"tech_stack"`
	Security      map[string]interface{} `json:"security"`
	Timestamp     time.Time              `json:"timestamp"`
}

func runWebChecks(ctx context.Context, target string) (*WebResult, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	// Add https:// if no scheme
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}

	// Parse URL to validate it
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}
	if parsedURL.Host == "" {
		return nil, fmt.Errorf("no host in URL")
	}

	result := &WebResult{
		URL:       target,
		Headers:   make(map[string]string),
		Redirects: []string{},
		TechStack: []string{},
		Security:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// Create HTTP client with timeout and redirect tracking
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Track redirects
			if len(via) > 0 {
				result.Redirects = append(result.Redirects, via[len(via)-1].URL.String())
			}
			// Limit redirects to 10
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	// Make the request
	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return nil, err
	}

	// Set common headers to avoid bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Store status
	result.StatusCode = resp.StatusCode
	result.StatusText = resp.Status

	// Store headers
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = strings.Join(values, ", ")
		}
	}

	// Check for security headers
	checkSecurityHeaders(result)

	// Detect technologies
	detectTechnologies(result, resp)

	// Read and analyze body content
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Read up to 1MB
	if err == nil {
		analyzeBodyContent(result, string(body))
	}

	return result, nil
}

func checkSecurityHeaders(result *WebResult) {
	securityHeaders := []string{
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Content-Security-Policy",
		"Referrer-Policy",
		"Permissions-Policy",
		"Expect-CT",
		"Public-Key-Pins",
		"Set-Cookie",
	}

	for _, header := range securityHeaders {
		if value, exists := result.Headers[header]; exists {
			result.Security[header] = value
		} else {
			result.Security[header] = "Missing"
		}
	}
}

func detectTechnologies(result *WebResult, resp *http.Response) {
	// Check server header
	if server, exists := result.Headers["Server"]; exists {
		result.TechStack = append(result.TechStack, "Server: "+server)
	}

	// Check X-Powered-By
	if poweredBy, exists := result.Headers["X-Powered-By"]; exists {
		result.TechStack = append(result.TechStack, "Powered by: "+poweredBy)
	}

	// Check Content-Type for hints
	if contentType, exists := result.Headers["Content-Type"]; exists {
		if strings.Contains(contentType, "PHP") {
			result.TechStack = append(result.TechStack, "PHP")
		}
	}
}

func analyzeBodyContent(result *WebResult, body string) {
	// Check for common meta tags
	if strings.Contains(body, "<meta name=\"generator\"") {
		// Extract generator
		start := strings.Index(body, "<meta name=\"generator\" content=\"")
		if start != -1 {
			start += len("<meta name=\"generator\" content=\"")
			end := strings.Index(body[start:], "\"")
			if end != -1 {
				generator := body[start : start+end]
				result.TechStack = append(result.TechStack, "Generator: "+generator)
			}
		}
	}

	// Check for common JavaScript frameworks
	javascriptFrameworks := []string{
		"react", "angular", "vue", "jquery", "bootstrap", "foundation",
	}
	
	for _, framework := range javascriptFrameworks {
		if strings.Contains(strings.ToLower(body), framework) {
			// Avoid false positives
			if !strings.Contains(strings.ToLower(body), "//") && 
			   !strings.Contains(strings.ToLower(body), "http") {
				result.TechStack = append(result.TechStack, framework)
			}
		}
	}
}

func formatWebResult(result *WebResult) string {
	var sb strings.Builder

	// Basic Info
	sb.WriteString(fmt.Sprintf("\u001b[0;36mURL:\u001b[0m %s\n", result.URL))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mStatus:\u001b[0m %d %s\n", result.StatusCode, result.StatusText))

	// Redirects
	if len(result.Redirects) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Redirects ===\u001b[0m\n")
		for i, redirect := range result.Redirects {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%d. %s\u001b[0m\n", i+1, redirect))
		}
	}

	// Headers
	if len(result.Headers) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Headers ===\u001b[0m\n")
		for key, value := range result.Headers {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%s:\u001b[0m %s\n", key, value))
		}
	}

	// Security
	if len(result.Security) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Security ===\u001b[0m\n")
		for key, value := range result.Security {
			status := "\u001b[0;32m[OK]\u001b[0m"
			if value == "Missing" {
				status = "\u001b[0;31m[MISSING]\u001b[0m"
			}
			sb.WriteString(fmt.Sprintf("%s \u001b[0;33m%s:\u001b[0m %v\n", status, key, value))
		}
	}

	// Tech Stack
	if len(result.TechStack) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Technology Stack ===\u001b[0m\n")
		for _, tech := range result.TechStack {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m- %s\u001b[0m\n", tech))
		}
	}

	return sb.String()
}
