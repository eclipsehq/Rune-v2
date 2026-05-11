package cmds

import (
	"fmt"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["snipe"] = Command{
		Category:    "information",
		Description: "Displays the last deleted message in the current channel.",
		Aliases:     []string{"s"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			CacheMu.Lock()
			deleted, ok := SnipeCache[m.ChannelID]
			CacheMu.Unlock()

			if !ok {
				msg.SendError(s, m, "Snipe", "There are no deleted messages to snipe in this channel.")
				return
			}

			content := deleted.Content
			if content == "" { content = "[No text content]" }

			output := fmt.Sprintf("User: %s#%s\nContent: %s", deleted.Author.Username, deleted.Author.Discriminator, content)
			msg.SendResponse(s, m, "Snipe", output)
		},
	}
}