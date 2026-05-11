package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"rune/internal/cmds"
	"rune/internal/auth"
	"rune/internal/config"
	"rune/internal/web"
	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := config.LoadConfig("cfg/config.json"); err != nil {
		log.Printf("Notice: Starting without existing config: %v", err)
	}

	var dg *discordgo.Session
	
	web.OnTokenChange = func(newToken string) {
		log.Println("[AUTH] Token change detected, re-authenticating...")
		if dg != nil {
			dg.Close()
		}
		
		var err error
		dg, err = auth.Authenticate(newToken)
		if err != nil {
			log.Printf("Error creating Discord session: %v", err)
			return
		}

		web.SetSession(dg)
		dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			fmt.Printf("Self-bot is now active. Logged in as: %s#%s (ID: %s)\n", r.User.Username, r.User.Discriminator, r.User.ID)
		})
		dg.AddHandler(cmds.TrackMessages)
		dg.AddHandler(cmds.HandleDelete)
		dg.AddHandler(cmds.HandleUpdate)
		dg.AddHandler(cmds.Handle)

		if err := dg.Open(); err != nil {
			log.Printf("Error opening connection: %v", err)
		} else {
			fmt.Println("Connection established.")
		}
	}

	web.Start()

	if config.Cfg.Token != "" && config.Cfg.Token != "YOUR_USER_TOKEN_HERE" {
		web.OnTokenChange(config.Cfg.Token)
	} else {
		log.Println("[!] No Discord token found. Configure via the dashboard.")
	}

	web.OpenDashboard()

	fmt.Println("Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	
	fmt.Println("Shutting down...")
	if dg != nil {
		dg.Close()
	}
}
