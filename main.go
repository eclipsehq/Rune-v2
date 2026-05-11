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
	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := config.LoadConfig("cfg/config.json"); err != nil {
		log.Fatalf("Critical error loading config: %v", err)
	}

	if config.Cfg.Token == "" || config.Cfg.Token == "YOUR_USER_TOKEN_HERE" {
		log.Fatal("Please provide a valid user token in config.json")
	}
	if config.Cfg.OwnerID == "" || config.Cfg.OwnerID == "YOUR_DISCORD_USER_ID" {
		log.Fatal("Please provide your Discord user ID as 'owner_id' in config.json. This is crucial for owner-only commands.")
	}
	if config.Cfg.Prefix == "" {
		log.Println("No prefix provided in config.json, defaulting to '!'")
		config.Cfg.Prefix = "!" 
	}
	dg, err := auth.Authenticate(config.Cfg.Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Self-bot is now active. Logged in as: %s#%s (ID: %s)\n", r.User.Username, r.User.Discriminator, r.User.ID)
		if r.User.ID != config.Cfg.OwnerID {
			log.Printf("WARNING: The bot's user ID (%s) does not match the configured owner ID (%s). Owner-only commands might not work as expected.", r.User.ID, config.Cfg.OwnerID)
		}
	})
	dg.AddHandler(cmds.TrackMessages)
	dg.AddHandler(cmds.HandleDelete)
	dg.AddHandler(cmds.HandleUpdate)
	dg.AddHandler(cmds.Handle)
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	fmt.Println("Connection established. The account will remain online until the program is closed.")
	fmt.Println("Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	
	fmt.Println("Shutting down...")
	dg.Close()
}
