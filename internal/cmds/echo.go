package cmds

import (
	"strings"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["echo"] = Command{
		Category:    "fun",
		Description: "Echoes back the provided text.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Echo", "Please provide a message to echo.")
				return
			}

			msg.SendMessage(s, m, strings.Join(args, " "))
		},
	}
}