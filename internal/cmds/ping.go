package cmds

import (
	"github.com/bwmarrin/discordgo"
	"rune/internal/msg"
)

func init() {
	Commands["ping"] = Command{
		Category:    "utility",
		Description: "Checks the bot's response time.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			msg.SendResponse(s, m, "Ping", "Pong!")
		},
	}
}