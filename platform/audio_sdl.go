// +build !null_audio

/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
)

type sdlAudioPlayer struct {
	music  map[string]*mix.Music
	sounds map[string]*mix.Chunk
}

func NewAudioPlayer() (AudioPlayer, error) {
	player := &sdlAudioPlayer{
		music:  make(map[string]*mix.Music),
		sounds: make(map[string]*mix.Chunk),
	}

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, mix.DEFAULT_CHUNKSIZE); err != nil {
		return player, err
	}

	return player, nil
}

func (player *sdlAudioPlayer) PlayMusic(name string, fade time.Duration, loops int) error {

	return nil
}

func (player *sdlAudioPlayer) LoadMusic(name string, data []byte) error {
	rwops := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	mus, err := mix.LoadMUSType_RW(rwops, mix.MID, 0)
	if err != nil {
		return err
	}

	player.music[name] = mus
	return nil
}

func (player *sdlAudioPlayer) LoadSound(name string, data []byte) error {
	rwops := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	chunk, err := mix.LoadWAV_RW(rwops, false)
	if err != nil {
		return err
	}

	player.sounds[name] = chunk
	return nil
}

func (player *sdlAudioPlayer) Shutdown() {
	for _, mus := range player.music {
		mus.Free()
	}

	for _, chunk := range player.sounds {
		chunk.Free()
	}

	mix.Quit()
}
