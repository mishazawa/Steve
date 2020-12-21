package main

import (
	"flag"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/google/uuid"
	"os"
	"strings"

	d "github.com/bwmarrin/discordgo"
)

var (
	token string
	host  string
	room  string
)

func init() {
	flag.StringVar(&token, "t", os.Getenv("TOKEN"), "Discord token")
	flag.StringVar(&host, "s", os.Getenv("HOST"), "Server host")
	flag.StringVar(&room, "r", os.Getenv("ROOM"), "Room id")
	flag.Parse()
}

type ChatTunnel struct {
	Discord   *d.Session
	Minecraft *bot.Client
}

func NewTunnel() *ChatTunnel {
	// create bots
	minecraft := bot.NewClient()
	discord, err := d.New("Bot " + token)

	if err != nil {
		panic(err)
	}

	return &ChatTunnel{discord, minecraft}
}

func (t *ChatTunnel) JoinServer(host string) error {
	return t.Minecraft.JoinServer(host, 25565)
}

func (t *ChatTunnel) JoinDiscord() error {
	t.Discord.Identify.Intents = d.MakeIntent(d.IntentsGuildMessages)

	return t.Discord.Open()
}

func (t *ChatTunnel) CloseDiscord() error {
	return t.Discord.Close()
}

func (t *ChatTunnel) HandleMessages(room string) error {
	// minecraft <-> discord
	t.Minecraft.Events.ChatMsg = onMinecraftChatMessage(t.Minecraft, t.Discord, room)
	t.Discord.AddHandler(onDiscordChatMessage(t.Minecraft, room))
	// handle game messages
	return t.Minecraft.HandleGame()
}

func main() {
	tunnel := NewTunnel()

	err := tunnel.JoinServer(host)
	if err != nil {
		panic(err)
	}

	err = tunnel.JoinDiscord()
	if err != nil {
		panic(err)
	}

	defer tunnel.CloseDiscord()

	err = tunnel.HandleMessages(room)
	if err != nil {
		panic(err)
	}
}

func onMinecraftChatMessage(minecraft *bot.Client, session *d.Session, channelId string) func(chat.Message, byte, uuid.UUID) error {
	return func(c chat.Message, pos byte, sender uuid.UUID) error {
		msg := c.ClearString()

		name := ""
		if user := strings.IndexByte(msg, '>'); user != -1 {
			name = msg[1:user]
		}
		if serv := strings.IndexByte(msg, ']'); serv != -1 {
			name = msg[1:serv]
		}
		content := msg
		if len(name) != 0 {
			content = msg[len(name)+3:]
		} else {
			name = "info"
		}

		if name == minecraft.Auth.Name {
			return nil
		}

		_, err := sendToDiscord(session, channelId, fmt.Sprintf("<%s> %s", name, content))

		return err
	}
}

func onDiscordChatMessage(minecraft *bot.Client, room string) func(*d.Session, *d.MessageCreate) {
	return func(s *d.Session, m *d.MessageCreate) {
		// skip self
		if m.Author.ID == s.State.User.ID {
			return
		}
		// skip non minecraft room
		if m.ChannelID != room {
			return
		}

		msg := chat.Text(fmt.Sprintf("@%s %s", m.Author.Username, m.Content))
		sendToMinecraft(minecraft, msg.String())
	}
}

func sendToMinecraft(minecraft *bot.Client, msg string) error {
	return minecraft.Chat(msg)
}

func sendToDiscord(session *d.Session, room, msg string) (*d.Message, error) {
	return session.ChannelMessageSend(room, msg)
}

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
