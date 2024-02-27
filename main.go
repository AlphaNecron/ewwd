package main

import (
	"ewwd/client"
	"ewwd/daemon"
	"os"
)

var handlers = map[string]struct {
	client func(args []string)
	daemon func()
}{
	"pa": {
		client: client.PulseAudio,
		daemon: daemon.PulseAudio,
	},
	"mp": {
		client: client.MusicPlayer,
		daemon: daemon.MusicPlayerDaemon,
	},
	"nm": {
		client: func(args []string) {

		},
		daemon: daemon.NetworkManager,
	},
}

func main() {
	if handler, ok := handlers[os.Args[1]]; ok {
		if len(os.Args) > 2 && os.Args[2] == "daemon" {
			handler.daemon()
		} else {
			handler.client(os.Args[2:])
		}
	}
}
