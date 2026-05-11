package cmds

import (
	"fmt"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["editsnipe"] = Command{
		Category:    "information",
		Description: "Displays the original content of the last edited message in the current channel.",
		Aliases:     []string{"es"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			CacheMu.Lock()
			edited, ok := EditSnipeCache[m.ChannelID]
			CacheMu.Unlock()

			if !ok {
				msg.SendError(s, m, "Edit Snipe", "There are no edited messages to snipe in this channel.")
				return
			}

			content := edited.Content
			if content == "" { content = "[No text content]" }

			output := fmt.Sprintf("User: %s#%s\nOriginal Content: %s", edited.Author.Username, edited.Author.Discriminator, content)
			msg.SendResponse(s, m, "Edit Snipe", output)
		},
	}
}