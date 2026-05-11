package cmds

import (
	"rune/internal/config"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["reload"] = Command{
		Category:    "utility",
		Description: "Reloads the config.json file.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if err := config.LoadConfig("cfg/config.json"); err != nil {
				msg.SendError(s, m, "Reload", "Failed to reload config: "+err.Error())
				return
			}
			msg.SendSuccess(s, m, "Reload", "Configuration successfully reloaded.")
		},
	}
}