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

package building

import (
	"image"
	"image/color"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

const buildingSize = 3 * 16

type HumanTownhall struct {
	Building

	worldPos           image.Point
	palette            color.Palette
	constructionSprite *resource.Sprite
	buildingSprite     *resource.Sprite
}

func NewHumanTownhall(res *resource.Resources, pal color.Palette) *HumanTownhall {
	bSpr := res.Sprites["HTOWNHAL.SPR"]
	coSpr := res.Sprites["HTHALLCO.SPR"]

	return &HumanTownhall{
		palette:            pal,
		constructionSprite: &coSpr,
		buildingSprite:     &bSpr,
	}
}

func (b *HumanTownhall) SetPosition(pt image.Point) {
	b.worldPos = pt
}

func (b *HumanTownhall) SetTilePosition(pt image.Point) {
	b.worldPos = pt.Mul(16)
}

func (b *HumanTownhall) Position() image.Point {
	return b.worldPos
}

func (b *HumanTownhall) Bounds() image.Rectangle {
	p := b.worldPos
	maxX := b.worldPos.X + buildingSize
	maxY := b.worldPos.Y + buildingSize
	return image.Rect(p.X, p.Y, maxX, maxY)
}

func (b *HumanTownhall) Update() error {
	return nil
}

func (b *HumanTownhall) Render(renderer platform.Renderer, cameraPos image.Point) error {
	img := b.buildingSprite.Imgs[0]
	renderer.BlitImage(b.worldPos.Add(cameraPos), img, b.palette)

	//g.renderer.BlitImage(g.cursorPos.Sub(cur.Point()), cur.Data, g.cursorPal)
	return nil
}
