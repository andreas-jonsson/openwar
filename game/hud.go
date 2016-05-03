/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"image"
	"image/color"
	"strings"

	"github.com/openwar-hq/openwar/resource"
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
		"IHRESBAR.IMG": {-72, 0},
		"IHRIGBAR.IMG": {-312, 0},
		"IHBOTBAR.IMG": {-72, -188},
		"IHLPANEL.IMG": {0, -72},
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
	hud.g.renderer.BlitPal(img.Data, hud.pal, gfx[name].Add(image.Point{img.X, img.Y}))
}
