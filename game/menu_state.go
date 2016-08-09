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

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
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
