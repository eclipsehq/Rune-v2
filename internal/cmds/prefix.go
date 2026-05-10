package cmds

import (
	"fmt"
	"log"
	"rune/internal/config"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["prefix"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		config.Mu.Lock()
		ownerID := config.Cfg.OwnerID
		currentPrefix := config.Cfg.Prefix
		config.Mu.Unlock()

		if m.Author.ID != ownerID {
			msg.SendResponse(s, m.ChannelID, "Prefix", "You are not authorized to change the prefix.")
			return
		}

		if len(args) == 0 {
			msg.SendResponse(s, m.ChannelID, "Prefix", fmt.Sprintf("Current prefix: %s. Usage: %sprefix <new>", currentPrefix, currentPrefix))
			return
		}

		newPrefix := args[0]
		if len(newPrefix) > 5 {
			msg.SendResponse(s, m.ChannelID, "Prefix", "Prefix cannot be longer than 5 characters.")
			return
		}

		config.Mu.Lock()
		config.Cfg.Prefix = newPrefix
		config.Mu.Unlock()

		msg.SendResponse(s, m.ChannelID, "Prefix", fmt.Sprintf("Prefix changed to: %s", newPrefix))
		log.Printf("Prefix changed to: %s by %s (%s)", newPrefix, m.Author.Username, m.Author.ID)
	}
}