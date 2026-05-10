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
	Commands["dicksize"] = Command{
		Category:    "fun",
		Description: "Generates a random dick size.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		size := rand.Intn(12) + 1
		shaft := ""
		for i := 0; i < size; i++ {
			shaft += "="
		}
		
		output := fmt.Sprintf("8%sD\nYour size: %d inches", shaft, size)
		msg.SendResponse(s, m, "Dick Size", output)
		},
	}
}