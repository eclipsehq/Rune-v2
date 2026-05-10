package cmds

import (
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["help"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		categories := map[string][]string{
			"utility":     {"ping", "prefix", "help"},
			"information": {},
			"fun":         {"echo"},
		}

		if len(args) == 0 {
			menu := "\u001b[38;5;33m[+]\u001b[0m FUN\n\u001b[38;5;33m[+]\u001b[0m UTILITY\n\u001b[38;5;33m[+]\u001b[0m INFORMATION"
			msg.SendResponse(s, m, "Help", menu)
			return
		}
		target := strings.ToLower(args[0])
		cmdList, exists := categories[target]

		if !exists {
			msg.SendResponse(s, m, "Help", " Category not found. Available: Utility, Information, FUN")
			return
		}

		header := strings.ToUpper(target)

		if len(cmdList) == 0 {
			msg.SendResponse(s, m, header, " No commands registered in this category.")
			return
		}

		output := " " + strings.Join(cmdList, ", ")
		msg.SendResponse(s, m, header, output)
	}
}