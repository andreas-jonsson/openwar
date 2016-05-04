/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"image"
	"time"

	"github.com/openwar-hq/openwar/resource"
)

type playState struct {
	g *Game
	p *player

	scroll float64

	ter *terrain
	res resource.Resources
}

func NewPlayState(g *Game) GameState {
	ter, _ := newTerrain(g)

	return &playState{
		g:   g,
		p:   newPlay(g, orcRace, g.resources.Palettes["FOREST.PAL"]),
		res: g.resources,
		ter: ter,
	}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from GameState, args ...interface{}) error {
	s.g.musicPlayer.random(1 * time.Second)

	snd, _ := s.g.soundPlayer.Sound("OREADY.VOC")
	snd.Play(-1, 0, 0)

	return nil
}

func (s *playState) Exit(to GameState) error {
	return nil
}

func (s *playState) Update() error {
	for {
		event := s.g.PollEvent()
		if event == nil {
			break
		}

		/*
			switch event.(type) {
			case *platform.QuitEvent:
				s.g.running = false
			}
		*/
	}
	return nil
}

func (s *playState) Render() error {
	s.scroll += s.g.dt * 0.005
	s.ter.render(image.Rect(0, 0, 320, 200), image.Point{int(s.scroll), 0})
	s.p.render()
	return nil
}
