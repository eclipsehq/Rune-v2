package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["slots"] = Command{
		Category:    "fun",
		Description: "Play the slot machine.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			emojis := []string{"🍎", "🍊", "🍐", "🍋", "🍉", "🍇", "🍓", "🍒"}
			res1 := emojis[rand.Intn(len(emojis))]
			res2 := emojis[rand.Intn(len(emojis))]
			res3 := emojis[rand.Intn(len(emojis))]

			result := fmt.Sprintf("[ %s | %s | %s ]", res1, res2, res3)
			
			status := "You lost!"
			if res1 == res2 && res2 == res3 {
				status = "JACKPOT!"
			} else if res1 == res2 || res2 == res3 || res1 == res3 {
				status = "So close!"
			}

			msg.SendResponse(s, m, "Slots", fmt.Sprintf("%s\n\n%s", result, status))
		},
	}
}