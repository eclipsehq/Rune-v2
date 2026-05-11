package cmds

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"rune/internal/msg"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["runefetch"] = Command{
		Category:    "information",
		Description: "Displays system information in a neofetch-style layout.",
		Aliases:     []string{"fetch", "rf"},
		Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			localIP := "127.0.0.1"
			addrs, err := net.InterfaceAddrs()
			if err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							localIP = ipnet.IP.String()
							break
						}
					}
				}
			}

			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			allocMB := mem.Alloc / 1024 / 1024
			totalMB := mem.Sys / 1024 / 1024

			cpuModel := getCPUModel()
			gpuModel := getGPUModel()

			uptime := time.Since(StartTime).Round(time.Second)

			logo := "\u001b[0;36m" +
				"  ____  _   _ _   _ _____\n" +
				" |  _ \\| | | | \\ | | ____|\n" +
				" | |_) | | | |  \\| |  _|\n" +
				" |  _ <| |_| | |\\  | |___\n" +
				" |_| \\_\\\\___/|_| \\_|_____|\u001b[0m\n"

			stats := fmt.Sprintf(
				"OS      :: %s (%s)\n"+
				"CPU     :: %s (%d Cores)\n"+
				"GPU     :: %s\n"+
				"Local IP:: %s\n"+
				"Memory  :: %dMB / %dMB\n"+
				"Uptime  :: %v\n"+
				"Runtime :: %s",
				runtime.GOOS, runtime.GOARCH, cpuModel, runtime.NumCPU(), gpuModel, localIP, allocMB, totalMB, uptime, runtime.Version(),
			)

			msg.SendResponse(s, m, "RuneFetch", logo+"\n"+stats)
		},
	}
}

func getCPUModel() string {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/cpuinfo")
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.Contains(line, "model name") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}
	return runtime.GOARCH
}

func getGPUModel() string {
	out, err := exec.Command("sh", "-c", "lspci | grep -i vga | cut -d: -f3").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	return "Unknown / Integrated"
}