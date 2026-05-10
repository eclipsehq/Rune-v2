package cmds

import (
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["flip"] = Command{
		Category:    "fun",
		Description: "Flips your text upside down.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "Flip", "Please provide text to flip.")
				return
			}

			normal := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.?!,;\"'()"
			flipped := "ɐqɔpǝɟƃɥᴉɾʞlɯuodbɹsʇnʌʍxʎz∀qƆpƎℲפHIſʞ˥WNOԀQᴚS┴∩ΛMX⅄Z0ƖᄅƐㄣϛ9ㄥ86˙¿¡'؛„,)( "
			
			input := strings.Join(args, " ")
			result := ""
			for _, char := range input {
				idx := strings.Index(normal, string(char))
				if idx != -1 {
					result = string(flipped[idx]) + result
				} else {
					result = string(char) + result
				}
			}
			msg.SendMessage(s, m, result)
		},
	}
}