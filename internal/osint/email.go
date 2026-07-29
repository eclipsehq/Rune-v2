package osint

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type EmailResult struct {
	Email        string                 `json:"email"`
	IsValid     bool                   `json:"is_valid"`
	Domain      string                 `json:"domain"`
	MXRecords   []string               `json:"mx_records"`
	SPFRecord   string                 `json:"spf_record"`
	DMARCRecord string                 `json:"dmarc_record"`
	TXTRecords  []string               `json:"txt_records"`
	Gravatar    string                 `json:"gravatar"`
	Timestamp   time.Time              `json:"timestamp"`
}

func runEmailChecks(ctx context.Context, target string) (*EmailResult, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	// Validate email format
	if !isValidEmail(target) {
		return &EmailResult{
			Email:    target,
			IsValid: false,
			Timestamp: time.Now(),
		}, fmt.Errorf("invalid email format")
	}

	result := &EmailResult{
		Email:    target,
		IsValid: true,
		Domain:  strings.Split(target, "@")[1],
		Timestamp: time.Now(),
	}

	// Get MX records
	if mxRecords, err := net.LookupMX(result.Domain); err == nil {
		for _, mx := range mxRecords {
			result.MXRecords = append(result.MXRecords, fmt.Sprintf("%d %s", mx.Pref, mx.Host))
		}
	}

	// Get TXT records (for SPF, DMARC, etc.)
	if txtRecords, err := net.LookupTXT(result.Domain); err == nil {
		for _, txtRecord := range txtRecords {
			// txtRecord is []string, join it
			joinedTxt := strings.Join(txtRecord, "")
			result.TXTRecords = append(result.TXTRecords, joinedTxt)
			
			// Check for SPF
			if strings.HasPrefix(strings.ToLower(joinedTxt), "v=spf1") {
				result.SPFRecord = joinedTxt
			}
			
			// Check for DMARC
			if strings.HasPrefix(strings.ToLower(joinedTxt), "v=dmarc1") {
				result.DMARCRecord = joinedTxt
			}
		}
	}

	// Check Gravatar
	result.Gravatar = getGravatarURL(target)

	return result, nil
}

func isValidEmail(email string) bool {
	// Simple email validation
	if !strings.Contains(email, "@") {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	local, domain := parts[0], parts[1]
	
	if local == "" || domain == "" {
		return false
	}
	
	// Check for valid domain format
	if !strings.Contains(domain, ".") {
		return false
	}
	
	// Check for spaces
	if strings.Contains(local, " ") || strings.Contains(domain, " ") {
		return false
	}
	
	return true
}

func getGravatarURL(email string) string {
	// Create MD5 hash of email (lowercase, trimmed)
	hash := md5Hash(strings.ToLower(strings.TrimSpace(email)))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s", hash)
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Check if SMTP server exists
func checkSMTPServer(ctx context.Context, domain string) (bool, error) {
	// Try to connect to SMTP server on port 25
	conn, err := smtp.Dial("mail." + domain + ":25")
	if err != nil {
		// Try without mail. prefix
		conn, err = smtp.Dial(domain + ":25")
		if err != nil {
			return false, nil
		}
	}
	defer conn.Quit()
	
	return true, nil
}

func formatEmailResult(result *EmailResult) string {
	var sb strings.Builder

	// Basic Info
	sb.WriteString(fmt.Sprintf("\u001b[0;36mEmail:\u001b[0m %s\n", result.Email))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mValid:\u001b[0m %v\n", result.IsValid))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mDomain:\u001b[0m %s\n", result.Domain))

	// MX Records
	if len(result.MXRecords) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Mail Servers (MX) ===\u001b[0m\n")
		for _, mx := range result.MXRecords {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m%s\u001b[0m\n", mx))
		}
	} else {
		sb.WriteString("\n\u001b[0;36m=== Mail Servers (MX) ===\u001b[0m\nNo MX records found\n")
	}

	// SPF Record
	if result.SPFRecord != "" {
		sb.WriteString("\n\u001b[0;36m=== SPF Record ===\u001b[0m\n")
		sb.WriteString(fmt.Sprintf("\u001b[0;33m%s\u001b[0m\n", result.SPFRecord))
	}

	// DMARC Record
	if result.DMARCRecord != "" {
		sb.WriteString("\n\u001b[0;36m=== DMARC Record ===\u001b[0m\n")
		sb.WriteString(fmt.Sprintf("\u001b[0;33m%s\u001b[0m\n", result.DMARCRecord))
	}

	// Gravatar
	if result.Gravatar != "" {
		sb.WriteString("\n\u001b[0;36m=== Gravatar ===\u001b[0m\n")
		sb.WriteString(fmt.Sprintf("\u001b[0;33mURL:\u001b[0m %s\n", result.Gravatar))
	}

	return sb.String()
}
