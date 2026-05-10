package auth

import (
	"github.com/bwmarrin/discordgo"
)


func Authenticate(token string) (*discordgo.Session, error) {
	dg, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	dg.Token = token
	return dg, nil
}
