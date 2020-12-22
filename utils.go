package main

import (
	"fmt"
	d "github.com/bwmarrin/discordgo"
)

// fn to find id of channel
func findDiscordChannelByName(s *d.Session, room string) (string, error) {
	guilds, err := s.UserGuilds(0, "", "")
	if err != nil {
		return "", err
	}
	for _, guild := range guilds {
		channels, err := s.GuildChannels(guild.ID)
		if err != nil {
			return "", err
		}
		for _, c := range channels {
			if c.Name == room {
				return c.ID, nil
			}
		}
	}
	return "", fmt.Errorf("Channel %s not found.", room)
}
