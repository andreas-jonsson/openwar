/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"math/rand"
	"path"
	"time"

	"github.com/openwar-hq/openwar/platform"
	"github.com/openwar-hq/openwar/resource"
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
