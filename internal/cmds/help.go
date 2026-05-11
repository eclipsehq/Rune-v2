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
			cats := make(map[string]int)
			for _, cmd := range Commands {
				cats[strings.ToUpper(cmd.Category)]++
			}
			var catList []string
			for c, count := range cats {
				catList = append(catList, fmt.Sprintf("\u001b[34m[+]\u001b[0m %-15s \u001b[0;30m(%d)\u001b[0m", c, count))
			}
			sort.Strings(catList)

			menu := strings.Join(catList, "\n")
			menu += "\n\n\u001b[0;34m[ INFO ]\u001b[0m"
			menu += fmt.Sprintf("\nTotal Commands :: %d", len(Commands))
			menu += fmt.Sprintf("\nUsage          :: %shelp [category]", prefix)
			msg.SendResponse(s, m, "Rune Help Menu", menu)
			return
		}

		query := strings.ToLower(args[0])

		var targetCmd *Command
		cmdName := query
		if cmd, ok := Commands[query]; ok {
			targetCmd = &cmd
		} else {
		findCmd:
			for name, cmd := range Commands {
				for _, alias := range cmd.Aliases {
					if strings.ToLower(alias) == query {
						targetCmd = &cmd
						cmdName = name
						break findCmd
					}
				}
			}
		}

		if targetCmd != nil {
			aliases := "None"
			if len(targetCmd.Aliases) > 0 {
				aliases = strings.Join(targetCmd.Aliases, ", ")
			}
			output := fmt.Sprintf("Command    :: %s%s\nCategory   :: %s\nAliases    :: %s\nDescription:: %s", 
				prefix, cmdName, strings.ToUpper(targetCmd.Category), aliases, targetCmd.Description)
			msg.SendResponse(s, m, "Command Info", output)
			return
		}

		var cmdList []string
		for name, cmd := range Commands {
			if strings.ToLower(cmd.Category) == query {
				entry := fmt.Sprintf("\u001b[34m[+]\u001b[0m %s%s - %s", prefix, name, cmd.Description)
				if len(cmd.Aliases) > 0 {
					entry += fmt.Sprintf(" \u001b[0;30m(%s)\u001b[0m", strings.Join(cmd.Aliases, ", "))
				}
				cmdList = append(cmdList, entry)
			}
		}
		sort.Strings(cmdList)

		if len(cmdList) == 0 {
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
		output := strings.Join(cmdList, "\n")
		output += "\n\n\u001b[0;34m[ INFO ]\u001b[0m"
		output += fmt.Sprintf("\nDetails        :: %shelp [command]", prefix)
		msg.SendResponse(s, m, "Category: "+header, output)
		},
	}
}