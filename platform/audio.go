/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import "time"

type AudioPlayer interface {
	LoadMusic(name string, data []byte) error
	PlayMusic(name string, fade time.Duration, loops int) error
	FadeMusic(fade time.Duration) error
	IsPlayingMusic() bool
	IsPausedMusic() bool
	StopMusic()
	PauseMusic()
	ResumeMusic()
	VolumeMusic(vol int) int
	MusicCallback(cb func())

	LoadSound(name string, data []byte) (Sound, error)
	Sound(name string) (Sound, error)

	StopChannel(channel int)
	IsPlayingChannel(channel int) bool
	FadeChannel(channel int, fade time.Duration) error
	VolumeChannel(channel, vol int) int
	ReserveChannels(num int) int
	ChannelCallback(cb func(int))

	Shutdown()
}

type Sound interface {
	String() string
	Length() time.Duration
	Play(channel, loops int, fade time.Duration) (int, error)
	Volume(vol int) int
}
