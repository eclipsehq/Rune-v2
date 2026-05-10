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
	Commands["calculate"] = Command{
		Category:    "utility",
		Description: "Evaluate a mathematical expression.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Calculate", "Please provide an expression.")
				return
			}

			expr := url.QueryEscape(strings.Join(args, " "))
			resp, err := http.Get("https://api.mathjs.org/v4/?expr=" + expr)
			if err != nil {
				msg.SendResponse(s, m, "Calculate", "Failed to evaluate expression.")
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			msg.SendResponse(s, m, "Calculate", string(body))
		},
	}
}