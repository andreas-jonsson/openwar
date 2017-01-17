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
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

func lerp(v0, v1, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}

func lerpRGBA(v0, v1 color.RGBA, t float64) (col color.RGBA) {
	col.R = uint8(lerp(float64(v0.R), float64(v1.R), t))
	col.G = uint8(lerp(float64(v0.G), float64(v1.G), t))
	col.B = uint8(lerp(float64(v0.B), float64(v1.B), t))
	col.A = uint8(lerp(float64(v0.A), float64(v1.A), t))
	return
}

type gradient image.Rectangle

func (g *gradient) ColorModel() color.Model {
	return color.RGBAModel
}

func (g *gradient) Bounds() image.Rectangle {
	return *(*image.Rectangle)(g)
}

func (g *gradient) At(x, y int) color.Color {
	top := color.RGBA{30, 30, 75, 255}
	bottom := color.RGBA{0, 0, 15, 255}

	t := float64(y) / float64(g.Max.Y)
	return lerpRGBA(top, bottom, t)
}

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
	renderer := s.g.renderer
	bb := renderer.BackBuffer()
	bounds := bb.Bounds()
	g := (*gradient)(&bounds)

	draw.Draw(bb, bounds, g, image.ZP, draw.Src)

	img := s.res.Images["TITLE.IMG"]
	pal := s.res.Palettes["TITLE.PAL"]

	offset := (bounds.Max.X - img.Data.Bounds().Max.X) / 2
	renderer.BlitImage(image.Point{offset, 0}, img.Data, pal)
	return nil
}
