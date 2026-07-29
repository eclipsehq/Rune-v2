# OSINT Framework Integration for Rune V2

This document describes the seamless integration of the OSINT Framework into Rune V2 Discord self-bot.

## Overview

The OSINT Framework has been integrated as a native module within Rune V2, providing powerful open-source intelligence gathering capabilities directly through Discord commands. This integration allows users to perform comprehensive reconnaissance and analysis on domains, IPs, emails, usernames, websites, and SSL certificates without leaving Discord.

## Features

### Available OSINT Commands

1. **`osint`** - Main OSINT command with subcommands
   - Usage: `osint [subcommand] [target]`
   - Subcommands: `domain`, `ip`, `email`, `username`, `web`, `ssl`

2. **`domain`** - Domain reconnaissance
   - Aliases: `dns`, `whois`, `rdap`
   - Usage: `domain [domain]` or `osint domain [domain]`
   - Features:
     - DNS record lookup (A, AAAA, MX, TXT, NS, CNAME, SOA, CAA)
     - Reverse DNS lookup
     - WHOIS information
     - RDAP information

3. **`ipinfo`** - IP address reconnaissance
   - Aliases: `iplookup`, `geoip`
   - Usage: `ipinfo [ip/hostname]` or `osint ip [ip/hostname]`
   - Features:
     - GeoIP information (country, region, city, coordinates)
     - ISP and organization details
     - ASN information
     - Reverse DNS lookup

4. **`email`** - Email address analysis
   - Usage: `email [email]` or `osint email [email]`
   - Features:
     - Email format validation
     - MX record lookup
     - SPF record detection
     - DMARC record detection
     - Gravatar URL generation

5. **`username`** - Username availability check
   - Aliases: `user`, `uname`, `checkuser`
   - Usage: `username [username]` or `osint username [username]`
   - Features:
     - Checks availability across 20+ platforms
     - Platforms: GitHub, GitLab, Twitter, Reddit, Instagram, Facebook, LinkedIn, YouTube, TikTok, Twitch, Discord, Steam, Spotify, Medium, Pinterest, Tumblr, Snapchat, WhatsApp, Telegram, Signal
     - Groups results by available/taken

6. **`web`** - Website analysis
   - Aliases: `website`, `site`, `url`
   - Usage: `web [url]` or `osint web [url]`
   - Features:
     - HTTP status code and headers
     - Redirect chain tracking
     - Security header analysis
     - Technology stack detection
     - Content analysis

7. **`ssl`** - SSL certificate analysis
   - Aliases: `cert`, `certificate`, `tls`
   - Usage: `ssl [hostname]` or `osint ssl [hostname]`
   - Features:
     - Certificate information (subject, issuer, serial number)
     - Validity period and expiration check
     - Subject Alternative Names (SANs)
     - TLS version detection
     - Cipher suite information
     - Self-signed certificate detection

## Command Usage Examples

```
# Domain reconnaissance
&domain example.com
&osint domain google.com

# IP lookup
&ipinfo 8.8.8.8
&osint ip github.com

# Email analysis
&email test@example.com
&osint email user@gmail.com

# Username availability
&username johndoe
&osint username alice123

# Website analysis
&web https://example.com
&osint web google.com

# SSL certificate analysis
&ssl example.com
&osint ssl github.com

# Get help on OSINT commands
&osint
```

## Configuration

The OSINT module uses default configuration values:

```go
OSINTConfig{
    Workers:   10,        // Concurrent workers for requests
    Timeout:   30,        // Timeout in seconds
    APIKeys:   map[string]string{}, // API keys for external services
}
```

### Adding API Keys

To enhance the OSINT capabilities, you can add API keys for external services in the configuration:

```go
cfg := osint.OSINTConfig{
    Workers: 10,
    Timeout: 30,
    APIKeys: map[string]string{
        "ipinfo": "your_ipinfo_api_key",
        "virustotal": "your_virustotal_api_key",
        "shodan": "your_shodan_api_key",
    },
}
```

## Cooldowns

To prevent abuse and rate limiting, commands have cooldowns:

- **General OSINT commands**: 5 seconds
- **Username check**: 10 seconds (makes multiple requests)

## Output Format

All OSINT commands return formatted output with ANSI color codes for better readability in Discord:

- **Blue**: Section headers
- **Yellow**: Field names
- **Green**: Available/positive results
- **Red**: Taken/negative results or errors

## Architecture

### Package Structure

```
internal/osint/
├── osint.go        # Main OSINT package with result types and configuration
├── domain.go       # Domain reconnaissance module
├── ip.go           # IP reconnaissance module
├── email.go        # Email analysis module
├── username.go     # Username availability module
├── web.go          # Web analysis module
└── ssl.go          # SSL certificate analysis module
```

### Integration Points

1. **Command Registration**: OSINT commands are registered in `internal/cmds/osint.go`
2. **Module Execution**: Each module runs in its own goroutine with context timeout
3. **Result Formatting**: Results are formatted with ANSI colors for Discord display
4. **Error Handling**: Errors are caught and displayed with appropriate messages

## Dependencies

The OSINT integration adds the following Go dependencies:

```
github.com/miekg/dns v1.1.59
```

This is used for DNS lookups in the domain module.

## Rate Limiting and Best Practices

1. **Use cooldowns**: Commands have built-in cooldowns to prevent abuse
2. **Respect rate limits**: External APIs may have rate limits
3. **Use responsibly**: Only perform OSINT on targets you have permission to investigate
4. **Legal compliance**: Ensure compliance with all applicable laws and service terms

## Future Enhancements

- [ ] Add more OSINT modules (VirusTotal, Shodan, Censys)
- [ ] Support for API key configuration via Discord commands
- [ ] Caching of results to reduce redundant requests
- [ ] Export results to files (JSON, CSV, Markdown)
- [ ] Interactive TUI for complex queries
- [ ] Batch processing of multiple targets

## Legal Disclaimer

This OSINT integration is intended for legal open-source intelligence gathering and security research on assets you own or have explicit permission to investigate. Users are responsible for complying with all applicable laws and service terms of use.

## License

The OSINT integration is part of Rune V2 and inherits its license.

---

**Integration Status**: ✅ Complete  
**Version**: 1.0.0  
**Last Updated**: 2024
