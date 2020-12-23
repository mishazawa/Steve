package utils

import (
	"fmt"
	d "github.com/bwmarrin/discordgo"
	"os"
	"path/filepath"
)

func GetPaths(root string) ([]string, error) {
	listDir := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			listDir = append(listDir, path)
		}
		return nil
	})

	return listDir, err
}

func ProcessList(arr []string, callback func(chan error, string)) error {
	errors := make(chan error)
	progress := 0

	for i := range arr {
		go callback(errors, arr[i])
	}

	for progress != len(arr) {
		select {
		case val := <-errors:
			if val != nil {
				return val
			}
			progress += 1
		}
	}
	return nil
}

// fn to find id of channel
func FindDiscordChannelByName(s *d.Session, room string) (string, error) {
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
