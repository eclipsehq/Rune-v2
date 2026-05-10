package cmds

import (
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["reverse"] = Command{
		Category:    "fun",
		Description: "Reverses the text you provide.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Reverse", "Please provide text to reverse.")
				return
			}

			input := strings.Join(args, " ")
			runes := []rune(input)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}

			msg.SendMessage(s, m, string(runes))
		},
	}
}