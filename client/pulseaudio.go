package client

import (
	"github.com/the-jonsey/pulseaudio"
	"log"
	"math"
	"strconv"
)

func PulseAudio(args []string) {
	// TODO: add check later
	vol := args[0]
	v, e := strconv.ParseFloat(vol, 64)
	if v > 100 || v < 0 {
		panic("Volume must be in range 0 and 1.")
	}
	if e != nil {
		panic(e)
	}
	client, e := pulseaudio.NewClient()
	if e != nil {
		panic(e)
	}
	defer client.Close()
	log.Fatalln(client.SetVolume(float32(math.Min(1, math.Ceil(v)/100))))
}
