package cmds

import (
	"rune/internal/config" // Ensure this folder exists at /internal/config/
	"strings"

	"github.com/bwmarrin/discordgo"
)


type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
var Commands = make(map[string]CommandFunc)

func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	config.Mu.Lock()
	prefix := config.Cfg.Prefix
	config.Mu.Unlock()

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	content := m.Content[len(prefix):]
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	name := strings.ToLower(parts[0])
	if fn, ok := Commands[name]; ok {
		fn(s, m, parts[1:])
	}
}