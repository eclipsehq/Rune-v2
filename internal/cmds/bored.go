package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type BoredResponse struct {
	Activity string `json:"activity"`
}

func init() {
	Commands["bored"] = Command{
		Category:    "fun",
		Description: "Suggests an activity when you are bored.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://www.boredapi.com/api/activity")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data BoredResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendResponse(s, m, "I'm Bored", data.Activity)
		},
	}
}