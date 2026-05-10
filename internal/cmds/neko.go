package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type NekoResponse struct {
	URL string `json:"neko"`
}

func init() {
	Commands["neko"] = Command{
		Category:    "fun",
		Description: "Fetches a random neko image.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://neko-love.xyz/api/v1/neko")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data NekoResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendMessage(s, m, data.URL)
		},
	}
}