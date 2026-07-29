package osint

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type DomainResult struct {
	DNS       map[string]interface{} `json:"dns"`
	WHOIS     map[string]interface{} `json:"whois"`
	RDAP      map[string]interface{} `json:"rdap"`
	Timestamp time.Time              `json:"timestamp"`
}

func runDomainChecks(ctx context.Context, target string) (*DomainResult, error) {
	// Clean the target
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	// Remove protocol prefixes if present
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimPrefix(target, "https://")
	target = strings.TrimPrefix(target, "www.")

	result := &DomainResult{
		DNS:       make(map[string]interface{}),
		WHOIS:     make(map[string]interface{}),
		RDAP:      make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// Run DNS checks
	if err := runDNSChecks(ctx, target, result); err != nil {
		// DNS checks failed, but we can continue with other checks
		result.DNS["error"] = err.Error()
	}

	// Run WHOIS check
	if whoisData, err := runWHOISCheck(ctx, target); err == nil {
		result.WHOIS = whoisData
	} else {
		result.WHOIS["error"] = err.Error()
	}

	// Run RDAP check
	if rdapData, err := runRDAPCheck(ctx, target); err == nil {
		result.RDAP = rdapData
	} else {
		result.RDAP["error"] = err.Error()
	}

	return result, nil
}

func runDNSChecks(ctx context.Context, target string, result *DomainResult) error {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	dnsData := make(map[string]interface{})

	// DNS record types to check
	recordTypes := []struct {
		name string
		dnsType uint16
	}{
		{"A", dns.TypeA},
		{"AAAA", dns.TypeAAAA},
		{"MX", dns.TypeMX},
		{"TXT", dns.TypeTXT},
		{"NS", dns.TypeNS},
		{"CNAME", dns.TypeCNAME},
		{"SOA", dns.TypeSOA},
		{"CAA", dns.TypeCAA},
	}

	for _, rt := range recordTypes {
		msg := new(dns.Msg)
		msg.SetQuestion(dns.Fqdn(target), rt.dnsType)
		
		r, _, err := c.ExchangeContext(ctx, msg, "8.8.8.8:53")
		if err != nil {
			continue
		}

		if len(r.Answer) > 0 {
			var records []string
			for _, ans := range r.Answer {
				switch rt.name {
				case "A":
					if a, ok := ans.(*dns.A); ok {
						records = append(records, a.A.String())
					}
				case "AAAA":
					if aaaa, ok := ans.(*dns.AAAA); ok {
						records = append(records, aaaa.AAAA.String())
					}
				case "MX":
					if mx, ok := ans.(*dns.MX); ok {
						records = append(records, fmt.Sprintf("%d %s", mx.Preference, mx.Mx))
					}
				case "TXT":
					if txt, ok := ans.(*dns.TXT); ok {
						records = append(records, strings.Join(txt.Txt, ""))
					}
				case "NS":
					if ns, ok := ans.(*dns.NS); ok {
						records = append(records, ns.Ns)
					}
				case "CNAME":
					if cname, ok := ans.(*dns.CNAME); ok {
						records = append(records, cname.Target)
					}
				case "SOA":
					if soa, ok := ans.(*dns.SOA); ok {
						records = append(records, fmt.Sprintf("%s %s %d %d %d %d %d", 
							soa.Ns, soa.Mx, soa.Serial, soa.Refresh, soa.Retry, soa.Expire, soa.Minimum))
					}
				case "CAA":
					if caa, ok := ans.(*dns.CAA); ok {
						records = append(records, fmt.Sprintf("%d %s %s", caa.Flag, caa.Tag, caa.Value))
					}
				}
			}
			if len(records) > 0 {
				dnsData[rt.name] = records
			}
		}
	}

	// Reverse DNS for first A record
	if aRecs, ok := dnsData["A"].([]string); ok && len(aRecs) > 0 {
		if names, err := net.LookupAddr(aRecs[0]); err == nil && len(names) > 0 {
			dnsData["reverse_dns"] = names
		}
	}

	result.DNS = dnsData
	return nil
}

func runWHOISCheck(ctx context.Context, target string) (map[string]interface{}, error) {
	// Simple WHOIS using net.LookupTXT as a fallback
	// In production, you'd use a proper WHOIS library
	
	// Try to get basic domain info
	ips, err := net.LookupIP(target)
	if err != nil {
		return nil, err
	}

	whoisData := make(map[string]interface{})
	if len(ips) > 0 {
		whoisData["ips"] = ips
	}

	// Add placeholder for full WHOIS implementation
	whoisData["note"] = "Full WHOIS lookup requires external library"
	
	return whoisData, nil
}

func runRDAPCheck(ctx context.Context, target string) (map[string]interface{}, error) {
	// RDAP lookup would go here
	// For now, return basic info
	
	rdapData := make(map[string]interface{})
	rdapData["note"] = "RDAP lookup requires external library"
	
	return rdapData, nil
}

func formatDomainResult(result *DomainResult) string {
	var sb strings.Builder

	// DNS Records
	sb.WriteString("\u001b[0;36m=== DNS Records ===\u001b[0m\n")
	for recordType, records := range result.DNS {
		if recordType == "error" {
			continue
		}
		sb.WriteString(fmt.Sprintf("\u001b[0;33m%s:\u001b[0m ", recordType))
		
		switch v := records.(type) {
		case []string:
			if len(v) > 0 {
				sb.WriteString(strings.Join(v, ", "))
			} else {
				sb.WriteString("No records found")
			}
		case string:
			sb.WriteString(v)
		default:
			sb.WriteString(fmt.Sprintf("%v", records))
		}
		sb.WriteString("\n")
	}

	// WHOIS Info
	if len(result.WHOIS) > 0 {
		sb.WriteString("\n\u001b[0;36m=== WHOIS Info ===\u001b[0m\n")
		for k, v := range result.WHOIS {
			if k == "error" {
				continue
			}
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%s:\u001b[0m %v\n", k, v))
		}
	}

	// RDAP Info
	if len(result.RDAP) > 0 {
		sb.WriteString("\n\u001b[0;36m=== RDAP Info ===\u001b[0m\n")
		for k, v := range result.RDAP {
			if k == "error" {
				continue
			}
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%s:\u001b[0m %v\n", k, v))
		}
	}

	return sb.String()
}
