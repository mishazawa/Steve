package main

import (
	"flag"
	"os"
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
