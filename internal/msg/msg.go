package msg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const Version = "2.1.0"

func send(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	_, err := s.ChannelMessageSendReply(m.ChannelID, content, &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, content)
	}
}

func Typing(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelTyping(m.ChannelID)
}

func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	send(s, m, content)
}

func SendResponse(s *discordgo.Session, m *discordgo.MessageCreate, funcName string, output string) {
	footer := fmt.Sprintf("\n\n\u001b[0;30mRune v%s\u001b[0m", Version)
	content := fmt.Sprintf("```ansi\n\u001b[0;36m[RUNE V2]\u001b[0m | %s\n```\n```ansi\n\u001b[0m%s%s\n```\n",
		funcName, output, footer)
	send(s, m, content)
}

func SendError(s *discordgo.Session, m *discordgo.MessageCreate, funcName string, output string) {
	footer := fmt.Sprintf("\n\n\u001b[0;30mRune v%s\u001b[0m", Version)
	content := fmt.Sprintf("```ansi\n\u001b[0;31m[ERROR]\u001b[0m | %s\n```\n```ansi\n\u001b[0m%s%s\n```\n",
		funcName, output, footer)
	send(s, m, content)
}

func SendSuccess(s *discordgo.Session, m *discordgo.MessageCreate, funcName string, output string) {
	footer := fmt.Sprintf("\n\n\u001b[0;30mRune v%s\u001b[0m", Version)
	content := fmt.Sprintf("```ansi\n\u001b[0;32m[SUCCESS]\u001b[0m | %s\n```\n```ansi\n\u001b[0m%s%s\n```\n",
		funcName, output, footer)
	send(s, m, content)
}
