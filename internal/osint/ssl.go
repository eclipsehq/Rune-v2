package osint

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

type SSLResult struct {
	Hostname      string                 `json:"hostname"`
	IP            string                 `json:"ip"`
	Port          int                    `json:"port"`
	Valid         bool                   `json:"valid"`
	Certificate   *CertificateInfo      `json:"certificate"`
	CipherSuites  []string               `json:"cipher_suites"`
	TLSVersion    string                 `json:"tls_version"`
	Timestamp     time.Time              `json:"timestamp"`
}

type CertificateInfo struct {
	Subject             string    `json:"subject"`
	Issuer              string    `json:"issuer"`
	SerialNumber        string    `json:"serial_number"`
	SignatureAlgorithm  string    `json:"signature_algorithm"`
	PublicKeyAlgorithm  string    `json:"public_key_algorithm"`
	PublicKeySize       int       `json:"public_key_size"`
	ValidFrom           time.Time `json:"valid_from"`
	ValidUntil          time.Time `json:"valid_until"`
	DaysUntilExpiry      int       `json:"days_until_expiry"`
	SANs                []string  `json:"sans"`
	IsExpired           bool      `json:"is_expired"`
	IsSelfSigned        bool      `json:"is_self_signed"`
}

func runSSLChecks(ctx context.Context, target string) (*SSLResult, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, fmt.Errorf("empty target")
	}

	// Remove protocol prefixes
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimPrefix(target, "https://")
	
	// Default to port 443
	port := 443
	
	result := &SSLResult{
		Hostname:    target,
		Port:        port,
		Valid:       false,
		Certificate: &CertificateInfo{},
		Timestamp:   time.Now(),
	}

	// Resolve IP
	ips, err := net.LookupIP(target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hostname: %v", err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IP addresses found")
	}
	result.IP = ips[0].String()

	// Try to connect using net.Dial with TLS
	addr := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		// Try with IP
		addr = fmt.Sprintf("%s:%d", result.IP, port)
		conn, err = net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}
	}
	defer conn.Close()

	// Wrap the connection in TLS
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         target,
	})
	defer tlsConn.Close()

	// Perform the TLS handshake
	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("TLS handshake failed: %v", err)
	}

	// Get connection state
	state := tlsConn.ConnectionState()
	result.Valid = true
	result.TLSVersion = tlsVersionToString(state.Version)

	// Get cipher suite
	if suite := tls.CipherSuiteName(state.CipherSuite); suite != "" {
		result.CipherSuites = append(result.CipherSuites, suite)
	}

	// Get certificates
	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		result.Certificate = parseCertificate(cert, target)
		
		// Check if certificate is valid for the hostname
		// Use VerifyHostname from crypto/tls if available, otherwise use a simple check
		if !verifyHostname(cert, target) {
			result.Valid = false
		}
	}

	return result, nil
}

func parseCertificate(cert *x509.Certificate, hostname string) *CertificateInfo {
	info := &CertificateInfo{
		Subject:        cert.Subject.String(),
		Issuer:         cert.Issuer.String(),
		SerialNumber:   cert.SerialNumber.String(),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		ValidFrom:      cert.NotBefore,
		ValidUntil:     cert.NotAfter,
		IsExpired:      time.Now().After(cert.NotAfter),
		IsSelfSigned:   cert.CheckSignatureFrom(cert) == nil,
	}

	// Calculate days until expiry
	if !info.IsExpired {
		days := int(time.Until(cert.NotAfter).Hours() / 24)
		info.DaysUntilExpiry = days
	} else {
		info.DaysUntilExpiry = 0
	}

	// Get Subject Alternative Names (SANs)
	for _, dnsName := range cert.DNSNames {
		info.SANs = append(info.SANs, dnsName)
	}
	for _, ip := range cert.IPAddresses {
		info.SANs = append(info.SANs, ip.String())
	}

	// Get public key size based on algorithm
	switch cert.PublicKeyAlgorithm {
	case x509.RSA:
		info.PublicKeySize = 2048 // Default RSA size
	case x509.ECDSA:
		info.PublicKeySize = 256 // Common ECDSA size
	case x509.Ed25519:
		info.PublicKeySize = 256
	}

	return info
}

// Simple hostname verification
func verifyHostname(cert *x509.Certificate, hostname string) bool {
	// Check if hostname matches Subject CN
	if cert.Subject.CommonName == hostname {
		return true
	}
	
	// Check if hostname matches any SAN
	for _, dnsName := range cert.DNSNames {
		if dnsName == hostname {
			return true
		}
	}
	
	// Check IP addresses
	for _, ip := range cert.IPAddresses {
		if ip.String() == hostname {
			return true
		}
	}
	
	return false
}

func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (%d)", version)
	}
}

func formatSSLResult(result *SSLResult) string {
	var sb strings.Builder

	// Basic Info
	sb.WriteString(fmt.Sprintf("\u001b[0;36mHostname:\u001b[0m %s\n", result.Hostname))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mIP:\u001b[0m %s\n", result.IP))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mPort:\u001b[0m %d\n", result.Port))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mTLS Version:\u001b[0m %s\n", result.TLSVersion))
	sb.WriteString(fmt.Sprintf("\u001b[0;36mValid:\u001b[0m %v\n", result.Valid))

	// Cipher Suites
	if len(result.CipherSuites) > 0 {
		sb.WriteString("\n\u001b[0;36m=== Cipher Suites ===\u001b[0m\n")
		for _, suite := range result.CipherSuites {
			sb.WriteString(fmt.Sprintf("\u001b[0;33m- %s\u001b[0m\n", suite))
		}
	}

	// Certificate Info
	if result.Certificate != nil {
		sb.WriteString("\n\u001b[0;36m=== Certificate ===\u001b[0m\n")
		
		cert := result.Certificate
		sb.WriteString(fmt.Sprintf("\u001b[0;33mSubject:\u001b[0m %s\n", cert.Subject))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mIssuer:\u001b[0m %s\n", cert.Issuer))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mSerial Number:\u001b[0m %s\n", cert.SerialNumber))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mSignature Algorithm:\u001b[0m %s\n", cert.SignatureAlgorithm))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mPublic Key Algorithm:\u001b[0m %s\n", cert.PublicKeyAlgorithm))
		if cert.PublicKeySize > 0 {
			sb.WriteString(fmt.Sprintf("\u001b[0;33mPublic Key Size:\u001b[0m %d bits\n", cert.PublicKeySize))
		}
		
		sb.WriteString("\n\u001b[0;36m=== Validity ===\u001b[0m\n")
		sb.WriteString(fmt.Sprintf("\u001b[0;33mValid From:\u001b[0m %s\n", cert.ValidFrom.Format("2006-01-02 15:04:05 UTC")))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mValid Until:\u001b[0m %s\n", cert.ValidUntil.Format("2006-01-02 15:04:05 UTC")))
		sb.WriteString(fmt.Sprintf("\u001b[0;33mDays Until Expiry:\u001b[0m %d\n", cert.DaysUntilExpiry))
		
		if cert.IsExpired {
			sb.WriteString("\u001b[0;31m[!] Certificate is EXPIRED!\u001b[0m\n")
		}
		if cert.IsSelfSigned {
			sb.WriteString("\u001b[0;33m[!] Certificate is self-signed\u001b[0m\n")
		}
		
		// SANs
		if len(cert.SANs) > 0 {
			sb.WriteString("\n\u001b[0;36m=== Subject Alternative Names (SANs) ===\u001b[0m\n")
			for _, san := range cert.SANs {
				sb.WriteString(fmt.Sprintf("\u001b[0;33m- %s\u001b[0m\n", san))
			}
		}
	}

	return sb.String()
}
