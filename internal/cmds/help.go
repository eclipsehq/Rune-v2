package cmds

import (
	"fmt"
	"rune/internal/config"
	"rune/internal/msg"
	"strings"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["help"] = Command{
		Category:    "utility",
		Description: "Displays the help menu.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		config.Mu.Lock()
		prefix := config.Cfg.Prefix
		config.Mu.Unlock()

		if len(args) == 0 {
			menu := fmt.Sprintf("Prefix: \u001b[38;5;33m%s\u001b[0m\n\n\u001b[38;5;33m[+]\u001b[0m FUN\n\u001b[38;5;33m[+]\u001b[0m UTILITY\n\u001b[38;5;33m[+]\u001b[0m INFORMATION", prefix)
			msg.SendResponse(s, m, "Help", menu)
			return
		}

		query := strings.ToLower(args[0])

		// 1. Check if it's a specific command lookup
		if cmd, ok := Commands[query]; ok {
			output := fmt.Sprintf("Command: %s%s\nCategory: %s\nDescription: %s", 
				prefix, query, strings.ToUpper(cmd.Category), cmd.Description)
			msg.SendResponse(s, m, "Command Info", output)
			return
		}

		// 2. Check if it's a category lookup
		var cmdList []string
		for name, cmd := range Commands {
			if strings.ToLower(cmd.Category) == query {
				cmdList = append(cmdList, fmt.Sprintf("%s%s - %s", prefix, name, cmd.Description))
			}
		}

		// Alphabetical sorting for a cleaner menu
		sort.Strings(cmdList)

		if len(cmdList) == 0 {
			// 3. Fallback: List available categories if the search failed
			cats := make(map[string]bool)
			for _, cmd := range Commands {
				cats[strings.ToUpper(cmd.Category)] = true
			}
			var catList []string
			for c := range cats {
				catList = append(catList, c)
			}
			sort.Strings(catList)

			msg.SendResponse(s, m, "Help", fmt.Sprintf("Command or Category not found.\nAvailable: %s", strings.Join(catList, ", ")))
			return
		}

		header := strings.ToUpper(query)
		output := " " + strings.Join(cmdList, "\n ")
		msg.SendResponse(s, m, header, output)
		},
	}
}