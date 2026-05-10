package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type CatResponse []struct {
	URL string `json:"url"`
}

func init() {
	Commands["cat"] = Command{
		Category:    "fun",
		Description: "Get a random cat image.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
			if err != nil {
				msg.SendResponse(s, m, "Cat", "Failed to fetch a cat image.")
				return
			}
			defer resp.Body.Close()

			var cat CatResponse
			json.NewDecoder(resp.Body).Decode(&cat)
			if len(cat) > 0 {
				msg.SendMessage(s, m, cat[0].URL)
			}
		},
	}
}