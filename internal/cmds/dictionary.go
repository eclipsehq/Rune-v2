package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DictionaryResponse []struct {
	Word     string `json:"word"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func init() {
	Commands["dictionary"] = Command{
		Category:    "information",
		Description: "Lookup the definition of a word.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Dictionary", "Please provide a word to look up.")
				return
			}

			word := args[0]
			resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
			if err != nil {
				msg.SendResponse(s, m, "Dictionary", "Failed to contact the dictionary API.")
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				msg.SendResponse(s, m, "Dictionary", "No definition found for that word.")
				return
			}

			var results DictionaryResponse
			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil || len(results) == 0 {
				msg.SendResponse(s, m, "Dictionary", "Failed to parse the definition.")
				return
			}

			res := results[0]
			var output strings.Builder
			output.WriteString(fmt.Sprintf("Word: %s\n", res.Word))

			for i, meaning := range res.Meanings {
				if i > 1 { break } 
				output.WriteString(fmt.Sprintf("\n(%s)\n", meaning.PartOfSpeech))
				if len(meaning.Definitions) > 0 {
					output.WriteString(fmt.Sprintf("- %s\n", meaning.Definitions[0].Definition))
				}
			}

			msg.SendResponse(s, m, "Dictionary", output.String())
		},
	}
}