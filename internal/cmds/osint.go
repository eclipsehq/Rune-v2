package cmds

import (
	"fmt"
	"rune/internal/msg"
	"rune/internal/osint"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	// OSINT Commands
	Commands["osint"] = Command{
		Category:    "osint",
		Description: "OSINT framework commands. Use subcommands: domain, ip, email, username, web, ssl.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "OSINT Framework", 
					"Available OSINT commands:\n"+
						"• osint domain [domain] - Domain reconnaissance\n"+
						"• osint ip [ip/hostname] - IP reconnaissance\n"+
						"• osint email [email] - Email analysis\n"+
						"• osint username [username] - Username availability\n"+
						"• osint web [url] - Web analysis\n"+
						"• osint ssl [hostname] - SSL certificate analysis")
				return
			}
			
			subcommand := strings.ToLower(args[0])
			
			// Remove the subcommand from args
			remainingArgs := args[1:]
			
			if len(remainingArgs) == 0 {
				msg.SendResponse(s, m, "OSINT "+strings.Title(subcommand), 
					fmt.Sprintf("Please provide a target for %s analysis.", subcommand))
				return
			}
			
			target := remainingArgs[0]
			
			// Use default OSINT config
			cfg := osint.DefaultConfig
			
			var result osint.Result
			
			switch subcommand {
			case "domain":
				result = osint.RunDomain(target, cfg)
			case "ip":
				result = osint.RunIP(target, cfg)
			case "email":
				result = osint.RunEmail(target, cfg)
			case "username":
				result = osint.RunUsername(target, cfg)
			case "web":
				result = osint.RunWeb(target, cfg)
			case "ssl":
				result = osint.RunSSL(target, cfg)
			default:
				msg.SendResponse(s, m, "OSINT", 
					fmt.Sprintf("Unknown OSINT subcommand: %s\nAvailable: domain, ip, email, username, web, ssl", subcommand))
				return
			}
			
			if result.Error != nil {
				msg.SendError(s, m, "OSINT "+strings.Title(subcommand), 
					fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			// Format the output for Discord
			output := fmt.Sprintf("Target: %s\nModule: %s\n\n%s", 
				result.Target, strings.ToUpper(result.Module), result.Output)
			
			msg.SendResponse(s, m, "OSINT "+strings.Title(subcommand), output)
		},
		Cooldown: 5 * time.Second, // 5 second cooldown to prevent abuse
	}

	// Individual OSINT commands for convenience
	Commands["domain"] = Command{
		Category:    "osint",
		Description: "Perform domain reconnaissance (DNS, WHOIS, RDAP).",
		Aliases:     []string{"dns", "whois", "rdap"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Domain Lookup", "Please provide a domain name.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunDomain(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "Domain Lookup", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Target: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "Domain Lookup", output)
		},
		Cooldown: 5 * time.Second,
	}

	Commands["ipinfo"] = Command{
		Category:    "osint",
		Description: "Get detailed information about an IP address or hostname.",
		Aliases:     []string{"iplookup", "geoip"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "IP Lookup", "Please provide an IP address or hostname.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunIP(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "IP Lookup", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Target: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "IP Lookup", output)
		},
		Cooldown: 5 * time.Second,
	}

	Commands["email"] = Command{
		Category:    "osint",
		Description: "Analyze an email address (MX, SPF, DMARC, Gravatar).",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Email Analysis", "Please provide an email address.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunEmail(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "Email Analysis", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Target: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "Email Analysis", output)
		},
		Cooldown: 5 * time.Second,
	}

	Commands["username"] = Command{
		Category:    "osint",
		Description: "Check username availability across multiple platforms.",
		Aliases:     []string{"user", "uname", "checkuser"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Username Check", "Please provide a username.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunUsername(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "Username Check", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Username: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "Username Check", output)
		},
		Cooldown: 10 * time.Second, // Longer cooldown as this makes many requests
	}

	Commands["web"] = Command{
		Category:    "osint",
		Description: "Analyze a website (headers, redirects, technology stack).",
		Aliases:     []string{"website", "site", "url"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Web Analysis", "Please provide a URL.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunWeb(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "Web Analysis", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Target: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "Web Analysis", output)
		},
		Cooldown: 5 * time.Second,
	}

	Commands["ssl"] = Command{
		Category:    "osint",
		Description: "Analyze SSL certificate of a website.",
		Aliases:     []string{"cert", "certificate", "tls"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "SSL Analysis", "Please provide a hostname.")
				return
			}
			
			cfg := osint.DefaultConfig
			result := osint.RunSSL(args[0], cfg)
			
			if result.Error != nil {
				msg.SendError(s, m, "SSL Analysis", fmt.Sprintf("Error: %v", result.Error))
				return
			}
			
			output := fmt.Sprintf("Target: %s\n\n%s", result.Target, result.Output)
			msg.SendResponse(s, m, "SSL Analysis", output)
		},
		Cooldown: 5 * time.Second,
	}
}
