package cmds

import (
	"fmt"
	"io"
	"net/http"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["weather"] = Command{
		Category:    "information",
		Description: "Checks the weather for a location.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Weather", "Please provide a location.")
				return
			}
			location := strings.Join(args, "+")
			resp, err := http.Get(fmt.Sprintf("https://wttr.in/%s?format=3", location))
			if err != nil {
				msg.SendResponse(s, m, "Weather", "Failed to get weather data.")
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			msg.SendResponse(s, m, "Weather", string(body))
		},
	}
}