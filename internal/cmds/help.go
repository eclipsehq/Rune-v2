package cmds

import (
	"fmt"
	"rune/internal/config"
	"rune/internal/msg"
	"strings"

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

		target := strings.ToLower(args[0])
		var cmdList []string

		for name, cmd := range Commands {
			if strings.ToLower(cmd.Category) == target {
				cmdList = append(cmdList, fmt.Sprintf("%s%s - %s", prefix, name, cmd.Description))
			}
		}

		if len(cmdList) == 0 {
			msg.SendResponse(s, m, "Help", " Category not found or empty. Available: Utility, Information, FUN")
			return
		}

		header := strings.ToUpper(target)
		output := " " + strings.Join(cmdList, "\n ")
		msg.SendResponse(s, m, header, output)
		},
	}
}