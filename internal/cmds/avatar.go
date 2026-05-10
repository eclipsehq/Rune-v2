package cmds

import (
	"rune/internal/msg"
	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["avatar"] = Command{
		Category:    "information",
		Description: "Displays the avatar of a user.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		target := m.Author
		if len(m.Mentions) > 0 {
			target = m.Mentions[0]
		} else if len(args) > 0 {
			u, err := s.User(args[0])
			if err == nil {
				target = u
			}
		}

		msg.SendResponse(s, m, "Avatar", target.AvatarURL("2048"))
		},
	}
}