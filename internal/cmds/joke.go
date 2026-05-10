package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type JokeResponse struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func init() {
	Commands["joke"] = Command{
		Category:    "fun",
		Description: "Get a random joke.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
			if err != nil {
				msg.SendResponse(s, m, "Joke", "Failed to fetch a joke.")
				return
			}
			defer resp.Body.Close()

			var joke JokeResponse
			json.NewDecoder(resp.Body).Decode(&joke)
			msg.SendResponse(s, m, "Joke", fmt.Sprintf("%s\n\n%s", joke.Setup, joke.Punchline))
		},
	}
}