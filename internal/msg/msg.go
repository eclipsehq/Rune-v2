package msg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)


func SendResponse(s *discordgo.Session, m *discordgo.MessageCreate, funcName string, output string) {
	content := fmt.Sprintf("```ansi\n\u001b[0;36m[RUNE V2]\u001b[0m | %s\n```\n```ansi\n\u001b[0m%s\n```\n",
		funcName, output)
	s.ChannelMessageSendReply(m.ChannelID, content, &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	})
}

func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	s.ChannelMessageSendReply(m.ChannelID, content, &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	})
}