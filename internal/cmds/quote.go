package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type QuoteResponse []struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

func init() {
	Commands["quote"] = Command{
		Category:    "fun",
		Description: "Gets a random inspirational quote.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://zenquotes.io/api/random")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data QuoteResponse
			json.NewDecoder(resp.Body).Decode(&data)
			if len(data) > 0 {
				msg.SendResponse(s, m, "Quote", fmt.Sprintf("\"%s\"\n- %s", data[0].Quote, data[0].Author))
			}
		},
	}
}