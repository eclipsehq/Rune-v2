package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["femboy"] = Command{
		Category:    "fun",
		Description: "Calculates your femboy percentage.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		percentage := rand.Intn(101)
		
		output := fmt.Sprintf("You are %d%% femboy.", percentage)
		msg.SendResponse(s, m, "Femboy Detector", output)
		},
	}
}