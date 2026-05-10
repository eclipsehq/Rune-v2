package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type DogResponse struct {
	URL string `json:"message"`
}

func init() {
	Commands["dog"] = Command{
		Category:    "fun",
		Description: "Fetches a random dog image.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data DogResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendMessage(s, m, data.URL)
		},
	}
}