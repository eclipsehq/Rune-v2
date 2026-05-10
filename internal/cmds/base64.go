package cmds

import (
	"encoding/base64"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["base64"] = Command{
		Category:    "utility",
		Description: "Encodes or decodes Base64 (Usage: &base64 encode/decode <text>).",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 2 {
				msg.SendResponse(s, m, "Base64", "Usage: &base64 <encode/decode> <text>")
				return
			}

			action := strings.ToLower(args[0])
			input := strings.Join(args[1:], " ")
			var result string

			if action == "encode" {
				result = base64.StdEncoding.EncodeToString([]byte(input))
			} else if action == "decode" {
				decoded, err := base64.StdEncoding.DecodeString(input)
				if err != nil {
					msg.SendResponse(s, m, "Base64 Error", "Invalid Base64 string.")
					return
				}
				result = string(decoded)
			} else {
				msg.SendResponse(s, m, "Base64", "Invalid action. Use 'encode' or 'decode'.")
				return
			}

			msg.SendResponse(s, m, "Base64 Result", result)
		},
	}
}