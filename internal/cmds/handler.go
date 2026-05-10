package cmds

import (
	"context"
	"rune/internal/config"
	"sync"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)


type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

type Command struct {
	Execute     CommandFunc
	Category    string
	Description string
	Aliases     []string
}

var Commands = make(map[string]Command)
var StartTime = time.Now()

var ActiveTasks = make(map[string]context.CancelFunc)
var TaskMu sync.Mutex

func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	config.Mu.Lock()
	prefix := config.Cfg.Prefix
	ownerID := config.Cfg.OwnerID
	config.Mu.Unlock()

	if m.Author.ID != ownerID || !strings.HasPrefix(m.Content, prefix) {
		return
	}

	content := m.Content[len(prefix):]
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	name := strings.ToLower(parts[0])
	if cmd, ok := Commands[name]; ok {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		cmd.Execute(s, m, parts[1:])
		return
	}

	for _, cmd := range Commands {
		for _, alias := range cmd.Aliases {
			if strings.ToLower(alias) == name {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				cmd.Execute(s, m, parts[1:])
				return
			}
		}
	}
}