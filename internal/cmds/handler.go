package cmds

import (
	"rune/internal/config"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)


type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

type Command struct {
	Execute     CommandFunc
	Category    string
	Description string
}

var Commands = make(map[string]Command)
var StartTime = time.Now()

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
	if cmd, ok := Commands[name]; ok {
		cmd.Execute(s, m, parts[1:])
	}
}