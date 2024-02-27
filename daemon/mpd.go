package daemon

import (
	"github.com/fhs/gompd/v2/mpd"
	"path/filepath"
	"strings"
	"time"
)

func PrintStatus(client *mpd.Client, callback func(attrs mpd.Attrs)) {
	status, e := client.Status()
	if e != nil {
		return
	}
	if callback != nil {
		callback(status)
	}
	song, e := client.CurrentSong()
	if e == nil {
		file := song["file"]
		status["currentsong"] = strings.TrimSuffix(file, filepath.Ext(file))
	}
	stdout.Encode(status)
}

func MusicPlayerDaemon() {
	client, e := mpd.Dial("tcp", ":6600")
	if e != nil {
		panic(e)
	}
	ticker := time.NewTicker(time.Second)
	keepAliveTicker := time.NewTicker(5 * time.Second)
	watcher, e := mpd.NewWatcher("tcp", ":6600", "", "player")
	isPaused := false
	if e != nil {
		panic(e)
	}
	callback := func(status mpd.Attrs) {
		pendingToSet := status["state"] != "play"
		if pendingToSet != isPaused {
			ticker.Reset(time.Second)
		}
		isPaused = pendingToSet
	}
	PrintStatus(client, callback)
	go func() {
		for range watcher.Event {
			PrintStatus(client, callback)
		}
	}()
	go func() {
		for range ticker.C {
			if !isPaused {
				PrintStatus(client, callback)
			}
		}
	}()
	go func() {
		for range keepAliveTicker.C {
			err := client.Ping()
			if err != nil {
				client.Close()
				keepAliveTicker.Stop()
				ticker.Stop()
				defer MusicPlayerDaemon()
				return
			}
		}
	}()
	select {}
}
