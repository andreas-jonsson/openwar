/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"image"
	"time"

	"github.com/openwar-hq/openwar/platform"
	"github.com/openwar-hq/openwar/resource"
)

type menuState struct {
	g   *Game
	res resource.Resources
}

func NewMenuState(g *Game) GameState {
	return &menuState{g: g, res: g.resources}
}

func (s *menuState) Name() string {
	return "menu"
}

func (s *menuState) Enter(from GameState, args ...interface{}) error {
	s.g.musicPlayer.play("MUSIC01.XMI", 3*time.Second)
	return nil
}

func (s *menuState) Exit(to GameState) error {
	s.g.musicPlayer.stop()
	return nil
}

func (s *menuState) Update() error {
	for {
		event := s.g.PollEvent()
		if event == nil {
			break
		}

		switch event.(type) {
		case *platform.KeyDownEvent, *platform.MouseButtonEvent:
			return s.g.SwitchState("play")
		}
	}
	return nil
}

func (s *menuState) Render() error {
	img := s.res.Images["TITLE.IMG"]
	pal := s.res.Palettes["TITLE.PAL"]

	s.g.renderer.BlitImage(image.Point{}, img.Data, pal)
	return nil
}
