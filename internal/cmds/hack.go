package cmds

import (
	"fmt"
	"math/rand"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["hack"] = Command{
		Category:    "fun",
		Description: "Perform a totally real hack on a user.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			target := "the server"
			if len(args) > 0 {
				target = args[0]
			}
			ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
			output := fmt.Sprintf("Bypassing firewall of %s...\nInjecting malicious SQL...\nFound credentials!\nIP: %s\nEmail: %s@gmail.com\nPassword: %s_123", 
				target, ip, target, target)
			msg.SendResponse(s, m, "Hacking...", output)
		},
	}
}