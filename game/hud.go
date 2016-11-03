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
	"image/draw"
	"strings"

	"github.com/andreas-jonsson/openwar/game/unit"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type (
	gameHudImpl struct {
		g      *Game
		race   playerRace
		pal    color.Palette
		images resource.Images
		units  *unit.Manager

		viewportBounds,
		miniMapViewportBounds image.Rectangle
		minimapFocus bool

		humanGfx, orcGfx map[string]image.Point
	}

	gameHud interface {
		handleMouse(mouse platform.MouseState) (image.Point, bool)
		handleMouseButton(ev *platform.MouseButtonEvent)
		render(miniMap *image.RGBA, cameraPos image.Point) error
		viewport() image.Rectangle
	}
)

func minimapRect() image.Rectangle {
	mmPos := image.Point{3, 6}
	bounds := image.Rect(0, 0, 64, 64).Add(mmPos)
	return bounds
}

func newGameHud(g *Game, units *unit.Manager, race playerRace, envPal color.Palette) gameHud {
	res := g.resources

	// Viewport is 240x176 in 4/3 mode and 336x176 in 16/9 mode.

	hud := &gameHudImpl{g: g, units: units, race: race}
	hud.viewportBounds = image.Rectangle{
		Min: image.Point{72, 12},
		Max: image.Point{312, 188},
	}

	addWidth := 0
	if g.config.Widescreen {
		maxX := 408
		addWidth = maxX - hud.viewportBounds.Max.X
		hud.viewportBounds.Max.X = maxX
	}

	hud.miniMapViewportBounds = image.Rectangle{image.ZP, hud.viewportBounds.Size().Div(16)}

	hud.images = make(resource.Images)
	hud.humanGfx = map[string]image.Point{
		"IHRESBAR.IMG": {72 + addWidth, 0},
		"IHRIGBAR.IMG": {312 + addWidth, 0},
		"IHBOTBAR.IMG": {72 + addWidth, 188},
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

func (hud *gameHudImpl) viewport() image.Rectangle {
	return hud.viewportBounds
}

func (hud *gameHudImpl) handleMouseButton(ev *platform.MouseButtonEvent) {
	bounds := minimapRect()

	if ev.Type == platform.MouseButtonDown && (image.Point{ev.X, ev.Y}).In(bounds) {
		hud.minimapFocus = true
	} else if ev.Type == platform.MouseButtonUp {
		hud.minimapFocus = false
	}
}

func (hud *gameHudImpl) handleMouse(mouse platform.MouseState) (image.Point, bool) {
	bounds := minimapRect()

	if hud.minimapFocus {
		if bounds.Min.X <= mouse.X && bounds.Max.X > mouse.X && bounds.Min.Y <= mouse.Y && bounds.Max.Y > mouse.Y {
			pos := image.Point{mouse.X, mouse.Y}
			pos = pos.Sub(bounds.Min)
			mmSize := hud.miniMapViewportBounds.Size().Div(2)
			return pos.Sub(mmSize).Mul(16), true
		}
	}

	return bounds.Min, false
}

func (hud *gameHudImpl) render(miniMap *image.RGBA, cameraPos image.Point) error {
	hud.patchWidescreen()
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

	mmPos := image.Point{3, 6}
	draw.Draw(hud.g.renderer.BackBuffer(), image.Rect(0, 0, 64, 64).Add(mmPos), miniMap, image.ZP, draw.Src)

	hud.renderMinimapUnits(mmPos)

	cameraPos = cameraPos.Div(16)
	hud.g.renderer.DrawRect(hud.miniMapViewportBounds.Add(cameraPos).Add(mmPos), color.RGBA{0x0, 0xFF, 0x0, 0xFF}, false)

	return nil
}

func (hud *gameHudImpl) renderMinimapUnits(mmPos image.Point) {
	for _, unit := range hud.units.AllUnits() {
		bounds := unit.Bounds()
		bounds.Min = bounds.Min.Div(16)
		bounds.Max = bounds.Max.Div(16)

		hud.g.renderer.DrawRect(bounds.Add(mmPos), color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, true)
	}
}

func (hud *gameHudImpl) renderImage(name string, gfx map[string]image.Point) {
	img := hud.images[name]
	hud.g.renderer.BlitImage(gfx[name], img.Data, hud.pal)

	//hud.g.renderer.DrawRect(hud.viewportBounds, color.RGBA{0xFF, 0, 0, 0xFF})
}

func (hud *gameHudImpl) patchWidescreen() {
	if hud.g.config.Widescreen {
		topPos := image.Point{72, 0}
		bottomPos := image.Point{72, 188}

		if hud.race == humanRace {
			hud.g.renderer.BlitImage(topPos, hud.images["IHRESBAR.IMG"].Data, hud.pal)
			hud.g.renderer.BlitImage(bottomPos, hud.images["IHBOTBAR.IMG"].Data, hud.pal)
		} else {
			hud.g.renderer.BlitImage(topPos, hud.images["IORESBAR.IMG"].Data, hud.pal)
			hud.g.renderer.BlitImage(bottomPos, hud.images["IOBOTBAR.IMG"].Data, hud.pal)
		}
	}
}
