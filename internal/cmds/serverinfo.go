package cmds

import (
	"fmt"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["serverinfo"] = Command{
		Category:    "information",
		Description: "Displays detailed information about the current server.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if m.GuildID == "" {
				msg.SendResponse(s, m, "Server Info", "This command can only be used within a server.")
				return
			}

			g, err := s.Guild(m.GuildID)
			if err != nil {
				msg.SendResponse(s, m, "Server Info", "Failed to retrieve server information.")
				return
			}

			creationTime, _ := discordgo.SnowflakeTimestamp(g.ID)

			output := fmt.Sprintf("Name: %s\nID: %s\nOwner ID: %s\nMembers: %d\nBoosts: %d (Tier %d)\nCreated: %s\nRegion: %s",
				g.Name, 
				g.ID, 
				g.OwnerID, 
				g.MemberCount, 
				g.PremiumSubscriptionCount, 
				int(g.PremiumTier), 
				creationTime.Format(time.RFC822), 
				g.Region)
			msg.SendResponse(s, m, "Server Info", output)
		},
	}
}