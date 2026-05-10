package msg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)


func SendResponse(s *discordgo.Session, channelID string, funcName string, output string) {
	msg := fmt.Sprintf("```ansi\n\u001b[0;36m[RUNE V2]\u001b[0m %s\n```\n```ansi\n\u001b[0m`%s\n```",
		funcName, output)
	s.ChannelMessageSend(channelID, msg)
}