/*
Copyright (C) 2016-2018 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
