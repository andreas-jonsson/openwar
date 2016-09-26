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
	"image/color"
)

type playerRace int

const (
	humanRace playerRace = iota
	orcRace
)

type player struct {
	g   *Game
	hud gameHud
}

func newPlay(g *Game, race playerRace, envPal color.Palette) *player {
	return &player{g: g, hud: newGameHud(g, race, envPal)}
}

func (p *player) render(miniMap *image.RGBA, cameraPos image.Point) error {
	return p.hud.render(miniMap, cameraPos)
}
