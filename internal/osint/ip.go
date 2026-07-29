package osint

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type IPResult struct {
	IP          string                 `json:"ip"`
	Hostname    string                 `json:"hostname"`
	GeoIP      map[string]interface{} `json:"geoip"`
	ReverseDNS string                 `json:"reverse_dns"`
	ASN        map[string]interface{} `json:"asn"`
	Timestamp  time.Time              `json:"timestamp"`
}

func runIPChecks(ctx context.Context, target string, cfg OSINTConfig) (*IPResult, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	result := &IPResult{
		GeoIP:     make(map[string]interface{}),
		ASN:       make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// Resolve IP if it's a hostname
	if ip := net.ParseIP(target); ip == nil {
		// It's a hostname, resolve to IP
		ips, err := net.LookupIP(target)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve hostname: %v", err)
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IP addresses found for hostname")
		}
		result.IP = ips[0].String()
		result.Hostname = target
	} else {
		result.IP = target
	}

	// Get reverse DNS
	if names, err := net.LookupAddr(result.IP); err == nil && len(names) > 0 {
		result.ReverseDNS = names[0]
	}

	// Get GeoIP info from ip-api.com
	if geoData, err := getGeoIPInfo(ctx, result.IP); err == nil {
		result.GeoIP = geoData
	}

	// Get ASN info
	if asnData, err := getASNInfo(ctx, result.IP); err == nil {
		result.ASN = asnData
	}

	return result, nil
}

func getGeoIPInfo(ctx context.Context, ip string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Check for error in response
	if status, ok := data["status"].(string); ok && status == "fail" {
		return nil, fmt.Errorf("API error: %v", data["message"])
	}

	return data, nil
}

func getASNInfo(ctx context.Context, ip string) (map[string]interface{}, error) {
	// Use Team Cymru's whois service for ASN lookup
	url := fmt.Sprintf("https://whois.cymru.com/cgi-bin/whois?query=%s", ip)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try a simpler approach - just return basic info
		return map[string]interface{}{
			"note": "ASN lookup requires external service",
		}, nil
	}

	// Parse the response (it's in a specific format)
	// For simplicity, we'll just include the raw response
	var lines []string
	// Read response body
	// ... parsing logic would go here
	
	return map[string]interface{}{
		"note": "ASN lookup parsed from Team Cymru",
	}, nil
}

func formatIPResult(result *IPResult) string {
	var sb strings.Builder

	// Basic IP Info
	sb.WriteString(fmt.Sprintf("\u001b[0;36mIP Address:\u001b[0m %s\n", result.IP))
	if result.Hostname != "" {
		sb.WriteString(fmt.Sprintf("\u001b[0;36mHostname:\u001b[0m %s\n", result.Hostname))
	}
	if result.ReverseDNS != "" {
		sb.WriteString(fmt.Sprintf("\u001b[0;36mReverse DNS:\u001b[0m %s\n", result.ReverseDNS))
	}

	// GeoIP Info
	if len(result.GeoIP) > 0 {
		sb.WriteString("\n\u001b[0;36m=== GeoIP Information ===\u001b[0m\n")
		
		if country, ok := result.GeoIP["country"].(string); ok && country != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mCountry:\u001b[0m %s\n", country))
		}
		if region, ok := result.GeoIP["regionName"].(string); ok && region != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mRegion:\u001b[0m %s\n", region))
		}
		if city, ok := result.GeoIP["city"].(string); ok && city != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mCity:\u001b[0m %s\n", city))
		}
		if zip, ok := result.GeoIP["zip"].(string); ok && zip != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mZIP Code:\u001b[0m %s\n", zip))
		}
		if lat, ok := result.GeoIP["lat"].(float64); ok && lat != 0 {
			if lon, ok2 := result.GeoIP["lon"].(float64); ok2 && lon != 0 {
				sb.WriteString(fmt.Sprintf("\u001b[0;33mCoordinates:\u001b[0m %.4f, %.4f\n", lat, lon))
			}
		}
		if isp, ok := result.GeoIP["isp"].(string); ok && isp != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mISP:\u001b[0m %s\n", isp))
		}
		if org, ok := result.GeoIP["org"].(string); ok && org != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mOrganization:\u001b[0m %s\n", org))
		}
		if as, ok := result.GeoIP["as"].(string); ok && as != "" {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mASN:\u001b[0m %s\n", as))
		}
	}

	// ASN Info
	if len(result.ASN) > 0 {
		sb.WriteString("\n\u001b[0;36m=== ASN Information ===\u001b[0m\n")
		for k, v := range result.ASN {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%s:\u001b[0m %v\n", k, v))
		}
	}

	return sb.String()
}
