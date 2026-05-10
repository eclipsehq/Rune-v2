package cmds

import (
	"fmt"
	"runtime"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["credits"] = Command{
		Category:    "information",
		Description: "Displays project credits and bot statistics.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			ramUsage := float64(mem.Alloc) / 1024 / 1024

			cmdCount := len(Commands)

			output := fmt.Sprintf("Developed with ❤️ by Light (eclipsehq) & Heavenzone\nA remaster of the eclipsehq legacy\n\nStats:\n- Commands: %d\n- RAM Usage: %.2f MB", 
				cmdCount, ramUsage)
			
			if len(args) > 0 && args[0] == "secret" {
				output += "\n\n||heavenzone the gawd||"
			}

			msg.SendResponse(s, m, "Credits & Stats", output)
		},
	}
}