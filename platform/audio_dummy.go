// +build js mobile

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

import (
	"errors"
	"time"
)

const MaxVolume = 0

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFadeMusic    = errors.New("can't fadeout music")
)

type dummySound struct{}

func (snd *dummySound) String() string {
	return ""
}

func (snd *dummySound) Length() time.Duration {
	return 0
}

func (snd *dummySound) Play(channel, loops int, fade time.Duration) (int, error) {
	return 0, nil
}

func (snd *dummySound) Volume(vol int) int {
	return 0
}

type dummyAudioPlayer struct{}

func NewAudioPlayer() (AudioPlayer, error) {
	return &dummyAudioPlayer{}, nil
}

func (player *dummyAudioPlayer) FadeMusic(fade time.Duration) error {
	return nil
}

func (player *dummyAudioPlayer) IsPlayingMusic() bool {
	return false
}

func (player *dummyAudioPlayer) IsPausedMusic() bool {
	return false
}

func (player *dummyAudioPlayer) StopMusic() {
}

func (player *dummyAudioPlayer) PauseMusic() {
}

func (player *dummyAudioPlayer) ResumeMusic() {
}

func (player *dummyAudioPlayer) VolumeMusic(vol int) int {
	return 0
}

func (player *dummyAudioPlayer) MusicCallback(cb func()) {
}

func (player *dummyAudioPlayer) PlayMusic(name string, fade time.Duration, loops int) error {
	return nil
}

func (player *dummyAudioPlayer) LoadMusic(name string, data []byte) error {
	return nil
}

func (player *dummyAudioPlayer) LoadSound(name string, data []byte) (Sound, error) {
	return &dummySound{}, nil
}

func (player *dummyAudioPlayer) Sound(name string) (Sound, error) {
	return &dummySound{}, nil
}

func (player *dummyAudioPlayer) StopChannel(channel int) {
}

func (player *dummyAudioPlayer) IsPlayingChannel(channel int) bool {
	return false
}

func (player *dummyAudioPlayer) FadeChannel(channel int, fade time.Duration) error {
	return nil
}

func (player *dummyAudioPlayer) VolumeChannel(channel, vol int) int {
	return 0
}

func (player *dummyAudioPlayer) ReserveChannels(num int) int {
	return 0
}

func (player *dummyAudioPlayer) ChannelCallback(cb func(int)) {
}

func (player *dummyAudioPlayer) Shutdown() {
}
