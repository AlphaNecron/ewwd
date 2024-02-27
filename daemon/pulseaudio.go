package daemon

import (
	"github.com/the-jonsey/pulseaudio"
	"time"
)

type audioState struct {
	Volume  float32 `json:"volume"`
	IsMuted bool    `json:"muted"`
}

var lastState = audioState{
	-1,
	false,
}

func printAudioStatus(client *pulseaudio.Client) {
	vol, e := client.Volume()
	if e != nil {
		return
	}
	muted, e := client.Mute()
	state := audioState{vol, muted}
	if lastState == state {
		return
	}
	lastState = state
	stdout.Encode(state)
}

func PulseAudio() {
	client, e := pulseaudio.NewClient()
	keepAliveTicker := time.NewTicker(time.Second * 5)
	if e != nil {
		panic(e)
	}
	c, e := client.UpdatesByType(pulseaudio.SUBSCRIPTION_MASK_SINK)
	if e != nil {
		panic(e)
	}
	printAudioStatus(client)
	go func() {
		for range keepAliveTicker.C {
			if !client.Connected() {
				client.Close()
				keepAliveTicker.Stop()
				defer PulseAudio()
				return
			}
		}
	}()
	for range c {
		printAudioStatus(client)
	}
}
