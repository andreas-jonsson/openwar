/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import "github.com/openwar-hq/openwar/resource"

type playState struct {
	g   *Game
	p   *player
	res resource.Resources
}

func NewPlayState(g *Game) GameState {
	return &playState{g: g, p: newPlay(g, orcRace, g.resources.Palettes["FOREST.PAL"]), res: g.resources}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from GameState, args ...interface{}) error {
	snd, _ := s.g.player.Sound("OREADY.VOC")
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
	s.p.render()
	return nil
}
