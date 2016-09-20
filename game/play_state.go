/*
Copyright (C) 2016 Andreas T Jonsson

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
	"image"
	"time"

	"github.com/andreas-jonsson/openwar/resource"
)

type playState struct {
	g *Game
	p *player

	scroll float64

	ter *terrain
	res resource.Resources
}

func NewPlayState(g *Game) GameState {
	ter, _ := newTerrain(g, "HUMAN01")

	return &playState{
		g:   g,
		p:   newPlay(g, humanRace, ter.pal),
		res: g.resources,
		ter: ter,
	}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from GameState, args ...interface{}) error {
	s.g.musicPlayer.random(10 * time.Second)

	snd, _ := s.g.soundPlayer.Sound("OREADY.VOC")
	snd.Play(-1, 0, 0)

	return nil
}

func (s *playState) Exit(to GameState) error {
	return nil
}

func (s *playState) Update() error {
	s.g.PollAll()
	return nil
}

func (s *playState) Render() error {
	s.scroll += s.g.dt * 0.005

	s.ter.render(s.g.renderer.BackBuffer().Bounds(), image.Point{int(s.scroll), 0})
	s.p.render()
	return nil
}
