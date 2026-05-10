package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["iq"] = Command{
		Category:    "fun",
		Description: "Calculates your IQ.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			iq := rand.Intn(201)
			msg.SendResponse(s, m, "IQ Test", fmt.Sprintf("Your IQ is: %d", iq))
		},
	}
}