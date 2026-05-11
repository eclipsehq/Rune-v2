package cmds

import (
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["status"] = Command{
		Category:    "utility",
		Description: "Changes the account's presence status (online, idle, dnd, invisible).",
		Aliases:     []string{"presence", "setstatus"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendError(s, m, "Status", "Usage: `status <online|idle|dnd|invisible>`")
				return
			}

			input := strings.ToLower(args[0])
			var status discordgo.Status

			switch input {
			case "online":
				status = discordgo.StatusOnline
			case "idle", "inactive", "away":
				status = discordgo.StatusIdle
			case "dnd", "do-not-disturb", "busy":
				status = discordgo.StatusDoNotDisturb
			case "invisible", "offline", "hidden":
				status = discordgo.StatusInvisible
			default:
				msg.SendError(s, m, "Status", "Invalid status. Options: `online`, `idle`, `dnd`, `invisible`.")
				return
			}

			_ = s.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status: string(status),
			})

			msg.SendSuccess(s, m, "Status", "Account status updated to **"+input+"**.")
		},
	}
}