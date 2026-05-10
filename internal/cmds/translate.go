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

func init() {
	Commands["translate"] = Command{
		Category:    "utility",
		Description: "Translates text to a specified language (e.g. &translate en hola).",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 2 {
				msg.SendResponse(s, m, "Translate", "Usage: translate <target_lang> <text>")
				return
			}

			targetLang := args[0]
			text := strings.Join(args[1:], " ")
			encodedText := url.QueryEscape(text)

			apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s", targetLang, encodedText)

			resp, err := http.Get(apiURL)
			if err != nil {
				msg.SendResponse(s, m, "Translate", "Failed to contact translation API.")
				return
			}
			defer resp.Body.Close()

			var result []interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				msg.SendResponse(s, m, "Translate", "Failed to parse translation response.")
				return
			}

			if len(result) > 0 {
				if inner, ok := result[0].([]interface{}); ok && len(inner) > 0 {
					if first, ok := inner[0].([]interface{}); ok && len(first) > 0 {
						if translatedText, ok := first[0].(string); ok {
							msg.SendResponse(s, m, "Translate", translatedText)
							return
						}
					}
				}
			}

			msg.SendResponse(s, m, "Translate", "Could not translate text.")
		},
	}
}