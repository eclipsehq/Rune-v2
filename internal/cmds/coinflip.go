package cmds

import (
	"math/rand"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	Commands["coinflip"] = Command{
		Category:    "fun",
		Description: "Flips a coin.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			result := "Heads"
			if rand.Intn(2) == 0 {
				result = "Tails"
			}

			msg.SendResponse(s, m, "Coin Flip", result)
		},
	}
}