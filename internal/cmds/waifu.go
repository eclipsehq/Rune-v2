package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type WaifuResponse struct {
	URL string `json:"url"`
}

func init() {
	Commands["waifu"] = Command{
		Category:    "fun",
		Description: "Get a random waifu image.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://api.waifu.pics/sfw/waifu")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data WaifuResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendMessage(s, m, data.URL)
		},
	}
}