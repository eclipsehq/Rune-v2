package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	Commands["ship"] = Command{
		Category:    "fun",
		Description: "Calculates the compatibility between two users.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 2 {
				msg.SendResponse(s, m, "Ship", "Please provide two names or mentions to ship.")
				return
			}

			name1 := args[0]
			name2 := args[1]
			percentage := rand.Intn(101)

			output := fmt.Sprintf("Shipping %s & %s\nCompatibility: %d%%", name1, name2, percentage)
			msg.SendResponse(s, m, "Ship", output)
		},
	}
}