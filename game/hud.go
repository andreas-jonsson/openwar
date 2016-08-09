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
	"strings"

	"github.com/andreas-jonsson/openwar/resource"
)

type gameHud struct {
	g      *Game
	race   playerRace
	pal    color.Palette
	images resource.Images

	humanGfx, orcGfx map[string]image.Point
}

func newGameHud(g *Game, race playerRace, envPal color.Palette) *gameHud {
	res := g.resources

	// Viewport is 240x176

	hud := &gameHud{g: g, race: race}
	hud.images = make(resource.Images)
	hud.humanGfx = map[string]image.Point{
		"IHRESBAR.IMG": {72, 0},
		"IHRIGBAR.IMG": {312, 0},
		"IHBOTBAR.IMG": {72, 188},
		"IHLPANEL.IMG": {0, 72},
		"IHMMAP01.IMG": {0, 0},
		"IHMMAP02.IMG": {0, 0},
	}

	if race == humanRace {
		hud.pal = resource.CombinePalettes(envPal, res.Palettes["SPRITE0.PAL"])
	} else {
		hud.pal = resource.CombinePalettes(envPal, res.Palettes["SPRITE1.PAL"])
	}
	hud.pal[0] = color.RGBA{}

	hud.orcGfx = make(map[string]image.Point, len(hud.humanGfx))
	for k, v := range hud.humanGfx {
		orcName := "IO" + strings.TrimPrefix(k, "IH")
		hud.orcGfx[orcName] = v

		hud.images[k] = res.Images[k]
		hud.images[orcName] = res.Images[orcName]
	}

	return hud
}

func (hud *gameHud) render() error {
	if hud.race == humanRace {
		hud.renderImage("IHRESBAR.IMG", hud.humanGfx)
		hud.renderImage("IHRIGBAR.IMG", hud.humanGfx)
		hud.renderImage("IHBOTBAR.IMG", hud.humanGfx)
		hud.renderImage("IHLPANEL.IMG", hud.humanGfx)
		hud.renderImage("IHMMAP01.IMG", hud.humanGfx)
	} else {
		hud.renderImage("IORESBAR.IMG", hud.orcGfx)
		hud.renderImage("IORIGBAR.IMG", hud.orcGfx)
		hud.renderImage("IOBOTBAR.IMG", hud.orcGfx)
		hud.renderImage("IOLPANEL.IMG", hud.orcGfx)
		hud.renderImage("IOMMAP01.IMG", hud.orcGfx)
	}
	return nil
}

func (hud *gameHud) renderImage(name string, gfx map[string]image.Point) {
	img := hud.images[name]
	hud.g.renderer.BlitImage(gfx[name].Add(image.Point{img.X, img.Y}), img.Data, hud.pal)
}
