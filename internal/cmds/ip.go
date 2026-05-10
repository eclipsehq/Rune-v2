package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rune/internal/msg"

	"github.com/bwmarrin/discordgo"
)

type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
	Message     string  `json:"message"`
}

func init() {
	Commands["ip"] = Command{
		Category:    "information",
		Description: "Retrieves information about an IP address.",
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		if len(args) == 0 {
			msg.SendResponse(s, m, "IP Lookup", "Please provide an IP address or domain.")
			return
		}

		resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", args[0]))
		if err != nil {
			msg.SendResponse(s, m, "IP Lookup", "Failed to reach the API.")
			return
		}
		defer resp.Body.Close()

		var info IPInfo
		json.NewDecoder(resp.Body).Decode(&info)

		if info.Status == "fail" {
			msg.SendResponse(s, m, "IP Lookup", fmt.Sprintf("Error: %s", info.Message))
			return
		}

		output := fmt.Sprintf("Query: %s\nLocation: %s, %s, %s\nISP: %s\nASN: %s", 
			info.Query, info.City, info.RegionName, info.Country, info.Isp, info.As)
		msg.SendResponse(s, m, "IP Lookup", output)
		},
	}
}