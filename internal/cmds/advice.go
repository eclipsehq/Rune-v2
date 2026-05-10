package cmds

import (
	"encoding/json"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type AdviceResponse struct {
	Slip struct {
		Advice string `json:"advice"`
	} `json:"slip"`
}

func init() {
	Commands["advice"] = Command{
		Category:    "fun",
		Description: "Get random life advice.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://api.adviceslip.com/advice")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var data AdviceResponse
			json.NewDecoder(resp.Body).Decode(&data)
			msg.SendResponse(s, m, "Advice", data.Slip.Advice)
		},
	}
}