package client

import (
	"github.com/fhs/gompd/v2/mpd"
	"log"
)

func throw(err error) {
	if err != nil {
		panic(err)
	}
}

func MusicPlayer(args []string) {
	client, e := mpd.Dial("tcp", ":6600")
	if e != nil {
		log.Fatalln(e)
	}
	defer client.Close()
	switch args[0] {
	case "play":
		throw(client.Play(-1))
		break
	case "pause":
		throw(client.Pause(true))
		break
	case "toggle":
		if status, e := client.Status(); e == nil {
			if status["state"] == "play" {
				throw(client.Pause(true))
			} else {
				throw(client.Play(-1))
			}
		}
		break
	}
}
