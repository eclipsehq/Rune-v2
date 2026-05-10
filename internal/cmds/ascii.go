package cmds

import (
	"io"
	"net/http"
	"net/url"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["ascii"] = Command{
		Category:    "fun",
		Description: "Convert text to ASCII art.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "ASCII", "Please provide text.")
				return
			}

			text := url.QueryEscape(strings.Join(args, " "))
			resp, err := http.Get("http://artii.herokuapp.com/make?text=" + text)
			if err != nil {
				msg.SendResponse(s, m, "ASCII", "Failed to generate ASCII art.")
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			msg.SendResponse(s, m, "ASCII", "\n"+string(body))
		},
	}
}