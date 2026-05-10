package cmds

import (
	"fmt"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["uptime"] = Command{
		Category:    "utility",
		Description: "Shows how long the bot has been running.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			uptime := time.Since(StartTime).Round(time.Second)
			output := fmt.Sprintf("Bot has been online for: %s", uptime.String())
			msg.SendResponse(s, m, "Uptime", output)
		},
	}
}