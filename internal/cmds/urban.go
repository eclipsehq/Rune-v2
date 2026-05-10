package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type UrbanResponse struct {
	List []struct {
		Definition string `json:"definition"`
		Example    string `json:"example"`
	} `json:"list"`
}

func init() {
	Commands["urban"] = Command{
		Category:    "information",
		Description: "Lookup a term on Urban Dictionary.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Urban", "Please provide a term to lookup.")
				return
			}

			term := url.QueryEscape(strings.Join(args, " "))
			resp, err := http.Get("http://api.urbandictionary.com/v0/define?term=" + term)
			if err != nil {
				msg.SendResponse(s, m, "Urban", "Failed to fetch definition.")
				return
			}
			defer resp.Body.Close()

			var urban UrbanResponse
			json.NewDecoder(resp.Body).Decode(&urban)
			if len(urban.List) == 0 {
				msg.SendResponse(s, m, "Urban", "No definition found.")
				return
			}

			def := strings.ReplaceAll(urban.List[0].Definition, "[", "")
			def = strings.ReplaceAll(def, "]", "")
			msg.SendResponse(s, m, "Urban", fmt.Sprintf("Definition:\n%s", def))
		},
	}
}