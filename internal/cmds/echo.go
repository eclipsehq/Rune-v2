package cmds

import (
	"strings"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["echo"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		if len(args) == 0 {
			msg.SendResponse(s, m.ChannelID, "Echo", "Please provide a message to echo.")
			return
		}

		msg.SendResponse(s, m.ChannelID, "Echo", strings.Join(args, " "))
	}
}