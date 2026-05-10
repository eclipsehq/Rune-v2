package cmds

import (
	"fmt"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["binary"] = Command{
		Category:    "utility",
		Description: "Converts text to binary.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Binary", "Please provide text.")
				return
			}
			var output strings.Builder
			for _, b := range []byte(strings.Join(args, " ")) {
				output.WriteString(fmt.Sprintf("%08b ", b))
			}
			msg.SendResponse(s, m, "Binary Output", output.String())
		},
	}
}