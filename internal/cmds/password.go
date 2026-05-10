package cmds

import (
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["password"] = Command{
		Category:    "utility",
		Description: "Generates a random 16-character password.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
			length := 16
			b := make([]byte, length)
			for i := range b {
				b[i] = chars[rand.Intn(len(chars))]
			}
			
			msg.SendResponse(s, m, "Password Generator", "I have sent a new password to your DMs.")
			dm, err := s.UserChannelCreate(m.Author.ID)
			if err == nil {
				s.ChannelMessageSend(dm.ID, "Generated Password: `"+string(b)+"`")
			}
		},
	}
}