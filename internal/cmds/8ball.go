package cmds

import (
	"math/rand"
	"rune/internal/msg"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	Commands["8ball"] = Command{
		Category:    "fun",
		Description: "Ask the magic 8-ball a question.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			responses := []string{
				"It is certain.", "It is decidedly so.", "Without a doubt.",
				"Yes definitely.", "You may rely on it.", "As I see it, yes.",
				"Most likely.", "Outlook good.", "Yes.", "Signs point to yes.",
				"Reply hazy, try again.", "Ask again later.", "Better not tell you now.",
				"Cannot predict now.", "Concentrate and ask again.", "Don't count on it.",
				"My reply is no.", "My sources say no.", "Outlook not so good.",
				"Very doubtful.",
			}

			if len(args) == 0 {
				msg.SendResponse(s, m, "8-Ball", "Please ask a question.")
				return
			}

			answer := responses[rand.Intn(len(responses))]
			msg.SendResponse(s, m, "8-Ball", answer)
		},
	}
}