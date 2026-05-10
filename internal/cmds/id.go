package cmds

import (
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["id"] = Command{
		Category:    "information",
		Description: "Gets the ID of a user, channel, or server.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			targetID := m.Author.ID
			label := "User ID"
			if len(m.Mentions) > 0 {
				targetID = m.Mentions[0].ID
			} else if len(args) > 0 && args[0] == "channel" {
				targetID = m.ChannelID
				label = "Channel ID"
			}
			msg.SendResponse(s, m, label, targetID)
		},
	}
}