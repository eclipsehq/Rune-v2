package cmds

import (
	"context"
	"fmt"
	"log"
	"rune/internal/config"
	"rune/internal/msg"
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
	Cooldown    time.Duration
}

var Commands = make(map[string]Command)
var StartTime = time.Now()

var cooldowns = make(map[string]time.Time)
var cdMu sync.Mutex

var ActiveTasks = make(map[string]context.CancelFunc)
var TaskMu sync.Mutex

var (
	MessageCache = make(map[string]*discordgo.Message)
	SnipeCache   = make(map[string]*discordgo.Message)
	EditSnipeCache = make(map[string]*discordgo.Message)
	CacheMu      sync.Mutex
)

func TrackMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	CacheMu.Lock()
	defer CacheMu.Unlock()
	MessageCache[m.ID] = m.Message
}

func HandleDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	CacheMu.Lock()
	defer CacheMu.Unlock()
	if msg, ok := MessageCache[m.ID]; ok {
		SnipeCache[m.ChannelID] = msg
		delete(MessageCache, m.ID)
	}
}

func HandleUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	CacheMu.Lock()
	defer CacheMu.Unlock()
	if oldMsg, ok := MessageCache[m.ID]; ok {
		if m.Content != "" && m.Content != oldMsg.Content {
			temp := *oldMsg
			EditSnipeCache[m.ChannelID] = &temp
			oldMsg.Content = m.Content
		}
	}
}

func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Centralized Recovery to prevent bot crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] Recovered in Handle: %v", r)
			msg.SendError(s, m, "System", "An internal error occurred while executing that command.")
		}
	}()

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
	
	var targetCmd *Command
	if cmd, ok := Commands[name]; ok {
		targetCmd = &cmd
	} else {
		for _, cmd := range Commands {
			for _, alias := range cmd.Aliases {
				if strings.ToLower(alias) == name {
					targetCmd = &cmd
					break
				}
			}
		}
	}

	if targetCmd != nil {
		// Cooldown Check
		cdKey := fmt.Sprintf("%s:%s", m.Author.ID, name)
		cdMu.Lock()
		if lastUsed, exists := cooldowns[cdKey]; exists && time.Since(lastUsed) < targetCmd.Cooldown {
			remaining := targetCmd.Cooldown - time.Since(lastUsed)
			cdMu.Unlock()
			msg.SendError(s, m, "Cooldown", fmt.Sprintf("Please wait %.1f seconds before using this again.", remaining.Seconds()))
			return
		}
		cooldowns[cdKey] = time.Now()
		cdMu.Unlock()

		s.ChannelMessageDelete(m.ChannelID, m.ID)
		targetCmd.Execute(s, m, parts[1:])
	}
}