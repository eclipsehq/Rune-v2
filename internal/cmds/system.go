package cmds

import (
	"fmt"
	"runtime"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["system"] = Command{
		Category:    "information",
		Description: "Displays system and OS information.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			output := fmt.Sprintf("OS: %s\nArch: %s\nCPU Cores: %d\nGoroutines: %d\nGo Version: %s",
				runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.NumGoroutine(), runtime.Version())
			
			msg.SendResponse(s, m, "System Info", output)
		},
	}
}