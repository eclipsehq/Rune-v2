package cmds

import (
	"context"
	"bufio"
	"fmt"
	"os"
	"rune/internal/msg"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["ap"] = Command{
		Category:    "fun",
		Description: "Auto chat-pack a user with words from words.txt for 30 seconds.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				msg.SendResponse(s, m, "AP", "Please provide a user to target.")
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			TaskMu.Lock()
			if oldCancel, ok := ActiveTasks["ap"]; ok { oldCancel() }
			ActiveTasks["ap"] = cancel
			TaskMu.Unlock()

			target := args[0]

			go func() {
				defer func() { TaskMu.Lock(); delete(ActiveTasks, "ap"); TaskMu.Unlock() }()
				var words []string
				file, err := os.Open("cfg/words.txt")
				if err != nil {
					words = []string{"WAKEUP"}
				} else {
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if line != "" {
							words = append(words, line)
						}
					}
					file.Close()
				}

				if len(words) == 0 {
					words = []string{"WAKEUP"}
				}

				endTime := time.Now().Add(30 * time.Second)
				i := 0
				for time.Now().Before(endTime) {
					select {
					case <-ctx.Done(): return
					default:
					word := words[i%len(words)]
					content := fmt.Sprintf("# %s %s", word, target)
					msg.SendMessage(s, m, content)
					i++
					time.Sleep(1200 * time.Millisecond)
					}
				}
			}()
		},
	}
}