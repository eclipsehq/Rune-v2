package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rune/internal/msg"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["crypto"] = Command{
		Category:    "information",
		Description: "Gets the current price of a cryptocurrency.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			coin := "bitcoin"
			if len(args) > 0 {
				coin = strings.ToLower(args[0])
			}
			resp, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin))
			if err != nil {
				msg.SendResponse(s, m, "Crypto", "API Error.")
				return
			}
			defer resp.Body.Close()
			var data map[string]map[string]float64
			json.NewDecoder(resp.Body).Decode(&data)
			if price, ok := data[coin]["usd"]; ok {
				msg.SendResponse(s, m, "Crypto", fmt.Sprintf("%s: $%.2f", strings.Title(coin), price))
			} else {
				msg.SendResponse(s, m, "Crypto", "Coin not found.")
			}
		},
	}
}