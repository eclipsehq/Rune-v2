package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type TrumpResponse struct {
	Value string `json:"value"`
}

func init() {
	Commands["trump"] = Command{
		Category:    "fun",
		Description: "Get a random Trump quote.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://api.tronalddump.io/random/quote")
			if err != nil {
				msg.SendResponse(s, m, "Trump Quote", "API Error.")
				return
			}
			defer resp.Body.Close()
			var data TrumpResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendResponse(s, m, "Trump says", data.Value)
		},
	}
}