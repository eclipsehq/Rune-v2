package cmds

import (
	"rune/internal/msg"
	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["stop"] = Command{
		Category:    "utility",
		Description: "Stops any running background tasks (like AP).",
		Aliases:     []string{"cancel", "halt"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			TaskMu.Lock()
			count := len(ActiveTasks)
			for name, cancel := range ActiveTasks {
				cancel()
				delete(ActiveTasks, name)
			}
			TaskMu.Unlock()

			if count > 0 {
				msg.SendSuccess(s, m, "Task Manager", "Successfully stopped all active tasks.")
			} else {
				msg.SendError(s, m, "Task Manager", "No active tasks found.")
			}
		},
	}
}