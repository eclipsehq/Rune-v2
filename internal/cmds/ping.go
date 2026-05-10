package cmds

import (
	"github.com/bwmarrin/discordgo"
	"rune/internal/msg"
)

func init() {
	Commands["ping"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		msg.SendResponse(s, m.ChannelID, "Ping", "Pong!")
	}
}