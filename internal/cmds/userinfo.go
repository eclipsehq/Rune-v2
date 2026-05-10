package cmds

import (
	"fmt"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["userinfo"] = Command{
		Category:    "information",
		Description: "Displays information about a user.",
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

		creationTime, _ := discordgo.SnowflakeTimestamp(target.ID)
		
		var joinedAt string = "N/A"
		if m.GuildID != "" {
			member, err := s.GuildMember(m.GuildID, target.ID)
			if err == nil && !member.JoinedAt.IsZero() {
				joinedAt = member.JoinedAt.Format(time.RFC822)
			}
		}

		output := fmt.Sprintf("Username: %s#%s\nID: %s\nCreated: %s\nJoined: %s\nBot: %t", 
			target.Username, target.Discriminator, target.ID, creationTime.Format(time.RFC822), joinedAt, target.Bot)
		msg.SendResponse(s, m, "User Info", output)
		},
	}
}