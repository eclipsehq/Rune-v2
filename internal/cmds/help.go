package cmds

import (
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["help"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		// Define command categories and their associated commands
		categories := map[string][]string{
			"utility":     {"ping", "prefix", "help"},
			"information": {},
			"fun":         {"echo"},
		}

		if len(args) == 0 {
			// Show main menu with 12-space indentation for sub-items
			// We lead with a space because SendResponse prepends a backtick
			menu := " Utility\n            Information\n            FUN"
			msg.SendResponse(s, m.ChannelID, "Help", menu)
			return
		}

		// Handle specific category request (e.g., &help utility)
		target := strings.ToLower(args[0])
		cmdList, exists := categories[target]

		if !exists {
			msg.SendResponse(s, m.ChannelID, "Help", " Category not found. Available: Utility, Information, FUN")
			return
		}

		header := strings.ToUpper(target)

		// Handle empty categories
		if len(cmdList) == 0 {
			msg.SendResponse(s, m.ChannelID, header, " No commands registered in this category.")
			return
		}

		// List commands for the selected category
		output := " " + strings.Join(cmdList, ", ")
		msg.SendResponse(s, m.ChannelID, header, output)
	}
}