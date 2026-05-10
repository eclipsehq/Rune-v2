package cmds

import (
	"strconv"
	"time"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["purge"] = Command{
		Category:    "utility",
		Description: "Deletes a specified number of your own messages in the current channel.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Purge", "Please specify the number of messages to scan.")
				return
			}

			amount, err := strconv.Atoi(args[0])
			if err != nil || amount <= 0 {
				msg.SendResponse(s, m, "Purge", "Please provide a valid positive number.")
				return
			}

			if amount > 100 {
				amount = 100
			}

			messages, err := s.ChannelMessages(m.ChannelID, amount, m.ID, "", "")
			if err != nil {
				return
			}

			for _, msgObj := range messages {
				if msgObj.Author.ID == s.State.User.ID {
					s.ChannelMessageDelete(m.ChannelID, msgObj.ID)
					time.Sleep(300 * time.Millisecond)
				}
			}

			s.ChannelMessageDelete(m.ChannelID, m.ID)
		},
	}
}