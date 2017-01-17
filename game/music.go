/*
Copyright (C) 2016-2017 Andreas T Jonsson

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

package game

import (
	"math/rand"
	"path"
	"time"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type musicPlayer struct {
	list     []string
	stopping bool
	player   platform.AudioPlayer
}

func newMusicPlayer(arch *resource.Archive, player platform.AudioPlayer) (*musicPlayer, error) {
	p := new(musicPlayer)
	p.list = make([]string, 0, 45)
	p.player = player

	for file := range arch.Files {
		if path.Ext(file) == ".XMI" {
			p.list = append(p.list, file)
		}
	}

	player.VolumeMusic(platform.MaxVolume / 2)
	return p, nil
}

func (p *musicPlayer) play(track string, fadein time.Duration) {
	p.player.PlayMusic(track, fadein, 0)
}

func (p *musicPlayer) random(fadein time.Duration) {
	p.player.PlayMusic(p.list[rand.Intn(len(p.list)-1)], fadein, 0)
	p.player.MusicCallback(func() {
		if !p.stopping {
			p.random(fadein)
		}
	})
}

func (p *musicPlayer) stop() {
	p.stopping = true
	p.player.StopMusic()
	p.stopping = false
}
