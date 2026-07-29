package osint

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type UsernameResult struct {
	Username     string                 `json:"username"`
	Platforms    map[string]bool       `json:"platforms"`
	Timestamp    time.Time              `json:"timestamp"`
}

// Platform availability check results
var platforms = []string{
	"github",
	"gitlab",
	"twitter",
	"reddit",
	"instagram",
	"facebook",
	"linkedin",
	"youtube",
	"tiktok",
	"twitch",
	"discord",
	"steam",
	"spotify",
	"medium",
	"pinterest",
	"tumblr",
	"snapchat",
	"whatsapp",
	"telegram",
	"signal",
}

func runUsernameChecks(ctx context.Context, target string) (*UsernameResult, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	result := &UsernameResult{
		Username:  target,
		Platforms: make(map[string]bool),
		Timestamp: time.Now(),
	}

	// Check availability on various platforms
	for _, platform := range platforms {
		available, err := checkUsernameAvailability(ctx, platform, target)
		if err == nil {
			result.Platforms[platform] = available
		} else {
			// If there's an error, we can't determine availability
			// For now, we'll mark it as unavailable
			result.Platforms[platform] = false
		}
	}

	return result, nil
}

func checkUsernameAvailability(ctx context.Context, platform, username string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	switch platform {
	case "github":
		url := fmt.Sprintf("https://github.com/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// If 404, username is available
		if resp.StatusCode == 404 {
			return true, nil
		}
		// If 200, username is taken
		if resp.StatusCode == 200 {
			return false, nil
		}
		// Other status codes - can't determine
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

	case "gitlab":
		url := fmt.Sprintf("https://gitlab.com/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

	case "twitter":
		// Twitter/X API - check if profile exists
		url := fmt.Sprintf("https://twitter.com/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// Twitter returns 404 for non-existent users
		if resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, nil

	case "reddit":
		url := fmt.Sprintf("https://www.reddit.com/user/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		// Add user agent to avoid bot detection
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// Reddit returns 404 for non-existent users
		if resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, nil

	case "instagram":
		url := fmt.Sprintf("https://www.instagram.com/%s/", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// Instagram returns 404 for non-existent users
		if resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, nil

	case "facebook":
		url := fmt.Sprintf("https://www.facebook.com/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// Facebook is tricky - might redirect or show a different page
		// For simplicity, we'll assume it's taken if we get 200
		if resp.StatusCode == 200 {
			return false, nil
		}
		return true, nil

	case "linkedin":
		url := fmt.Sprintf("https://www.linkedin.com/in/%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// LinkedIn returns 999 for non-existent profiles
		if resp.StatusCode == 999 || resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, nil

	case "youtube":
		url := fmt.Sprintf("https://www.youtube.com/@%s", username)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		
		// YouTube returns 404 for non-existent channels
		if resp.StatusCode == 404 {
			return true, nil
		}
		if resp.StatusCode == 200 {
			return false, nil
		}
		return false, nil

	default:
		// For platforms we don't have specific checks for,
		// we'll assume the username is available
		return true, nil
	}
}

func formatUsernameResult(result *UsernameResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\u001b[0;36mUsername:\u001b[0m %s\n", result.Username))
	sb.WriteString("\n")

	// Group platforms by availability
	var available []string
	var taken []string

	for platform, isAvailable := range result.Platforms {
		if isAvailable {
			available = append(available, platform)
		} else {
			taken = append(taken, platform)
		}
	}

	// Sort the lists
	sort.Strings(available)
	sort.Strings(taken)

	// Available platforms
	if len(available) > 0 {
		sb.WriteString("\u001b[0;32m[+] Available on:\u001b[0m\n")
		for _, platform := range available {
			sb.WriteString(fmt.Sprintf("  - %s\n", platform))
		}
		sb.WriteString("\n")
	}

	// Taken platforms
	if len(taken) > 0 {
		sb.WriteString("\u001b[0;31m[-] Taken on:\u001b[0m\n")
		for _, platform := range taken {
			sb.WriteString(fmt.Sprintf("  - %s\n", platform))
		}
	}

	return sb.String()
}
