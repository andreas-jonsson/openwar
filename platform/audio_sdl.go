/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"errors"
	"os"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
)

const MaxVolume = mix.MAX_VOLUME

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFadeMusic    = errors.New("can't fadeout music")
)

type sdlSound struct {
	name  string
	chunk *mix.Chunk
}

func (snd *sdlSound) String() string {
	return snd.name
}

func (snd *sdlSound) Length() time.Duration {
	return time.Duration(snd.chunk.LengthInMs()) * time.Millisecond
}

func (snd *sdlSound) Play(channel, loops int, fade time.Duration) (int, error) {
	if fade > 0 {
		return snd.chunk.FadeIn(channel, loops, int(fade/time.Millisecond))
	}
	return snd.chunk.Play(channel, loops)
}

func (snd *sdlSound) Volume(vol int) int {
	return snd.chunk.Volume(vol)
}

type sdlAudioPlayer struct {
	music  map[string]*mix.Music
	sounds map[string]*sdlSound
}

func NewAudioPlayer() (AudioPlayer, error) {
	player := &sdlAudioPlayer{
		music:  make(map[string]*mix.Music),
		sounds: make(map[string]*sdlSound),
	}

	wd, _ := os.Getwd()
	os.Setenv("TIMIDITY_CFG", wd+"data/timidity/timidity.cfg")

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, mix.DEFAULT_CHUNKSIZE); err != nil {
		return player, err
	}

	return player, nil
}

func (player *sdlAudioPlayer) FadeMusic(fade time.Duration) error {
	if !mix.FadeOutMusic(int(fade / time.Millisecond)) {
		return ErrFadeMusic
	}
	return nil
}

func (player *sdlAudioPlayer) IsPlayingMusic() bool {
	return mix.PlayingMusic()
}

func (player *sdlAudioPlayer) IsPausedMusic() bool {
	return mix.PausedMusic()
}

func (player *sdlAudioPlayer) StopMusic() {
	mix.HaltMusic()
}

func (player *sdlAudioPlayer) PauseMusic() {
	mix.PauseMusic()
}

func (player *sdlAudioPlayer) ResumeMusic() {
	mix.ResumeMusic()
}

func (player *sdlAudioPlayer) VolumeMusic(vol int) int {
	return mix.VolumeMusic(vol)
}

func (player *sdlAudioPlayer) MusicCallback(cb func()) {
	mix.HookMusicFinished(cb)
}

func (player *sdlAudioPlayer) PlayMusic(name string, fade time.Duration, loops int) error {
	mus, ok := player.music[name]
	if !ok {
		return ErrFileNotFound
	}

	var err error
	if fade > 0 {
		err = mus.FadeIn(loops, int(fade/time.Millisecond))
	} else {
		err = mus.Play(loops)
	}

	return err
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

func (player *sdlAudioPlayer) LoadSound(name string, data []byte) (Sound, error) {
	rwops := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	chunk, err := mix.LoadWAV_RW(rwops, false)
	if err != nil {
		return nil, err
	}

	snd := &sdlSound{name: name, chunk: chunk}
	player.sounds[name] = snd
	return snd, nil
}

func (player *sdlAudioPlayer) Sound(name string) (Sound, error) {
	snd, ok := player.sounds[name]
	if ok {
		return snd, nil
	}
	return snd, ErrFileNotFound
}

func (player *sdlAudioPlayer) StopChannel(channel int) {
	mix.HaltChannel(channel)
}

func (player *sdlAudioPlayer) IsPlayingChannel(channel int) bool {
	return mix.Playing(channel) != 0
}

func (player *sdlAudioPlayer) FadeChannel(channel int, fade time.Duration) error {
	mix.FadeOutChannel(channel, int(fade/time.Millisecond))
	return nil
}

func (player *sdlAudioPlayer) VolumeChannel(channel, vol int) int {
	return mix.Volume(channel, vol)
}

func (player *sdlAudioPlayer) ReserveChannels(num int) int {
	return mix.ReserveChannels(num)
}

func (player *sdlAudioPlayer) ChannelCallback(cb func(int)) {
	mix.ChannelFinished(cb)
}

func (player *sdlAudioPlayer) Shutdown() {
	for _, mus := range player.music {
		mus.Free()
	}

	for _, snd := range player.sounds {
		snd.chunk.Free()
	}

	mix.CloseAudio()
}
