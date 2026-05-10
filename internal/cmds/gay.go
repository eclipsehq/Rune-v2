package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["gay"] = Command{
		Category:    "fun",
		Description: "Calculates your gay percentage.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			percentage := rand.Intn(101)
			msg.SendResponse(s, m, "Gay Detector", fmt.Sprintf("You are %d%% gay.", percentage))
		},
	}
}