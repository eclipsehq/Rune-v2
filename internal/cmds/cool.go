package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["cool"] = Command{
		Category:    "fun",
		Description: "Calculates how cool you are.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			percentage := rand.Intn(101)
			msg.SendResponse(s, m, "Cool Detector", fmt.Sprintf("You are %d%% cool.", percentage))
		},
	}
}