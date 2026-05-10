package cmds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"rune/internal/msg"
)

func init() {
	Commands["ping"] = Command{
		Category:    "utility",
		Description: "Checks the bot's heartbeat latency.",
		Aliases:     []string{"latency", "ms"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			latency := s.HeartbeatLatency()
			msg.SendResponse(s, m, "Latency", fmt.Sprintf("Pong! 🏓\nHeartbeat: %v", latency))
		},
	}
}