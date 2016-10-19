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

package unit

import (
	"image"
	"image/color"
	"log"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type (
	Building struct {
		constructionSprite *resource.Sprite
		buildingSprite     *resource.Sprite

		config   *BuildingConfig
		id       uint64
		worldPos image.Point
		palette  color.Palette

		UnitName     string  `network:"sync"`
		Health       int     `network:"sync"`
		Construction float32 `network:"sync"`
	}

	BuildingConfig struct {
		Size int
		BuildingSprite,
		ConstructionSprite string
	}
)

func NewBuilding(name string, cfg *BuildingConfig, res *resource.Resources, pal color.Palette) *Building {
	bSpr, ok := res.Sprites[cfg.BuildingSprite]
	if !ok {
		log.Panicln("could not load:", cfg.BuildingSprite)
	}

	coSpr, ok := res.Sprites[cfg.ConstructionSprite]
	if !ok {
		log.Panicln("could not load:", cfg.ConstructionSprite)
	}

	return &Building{
		id:                 platform.NewId64(),
		config:             cfg,
		palette:            pal,
		constructionSprite: &coSpr,
		buildingSprite:     &bSpr,
		UnitName:           name,
	}
}

func (b *Building) Id() uint64 {
	return b.id
}

func (b *Building) Type() UnitType {
	return BuildingType
}

func (b *Building) Name() string {
	return b.UnitName
}

func (b *Building) SetPosition(pt image.Point) {
	b.worldPos = pt
}

func (b *Building) SetTilePosition(pt image.Point) {
	b.worldPos = pt.Mul(16)
}

func (b *Building) Position() image.Point {
	return b.worldPos
}

func (b *Building) Bounds() image.Rectangle {
	p := b.worldPos
	sz := b.config.Size * 16
	maxX := b.worldPos.X + sz
	maxY := b.worldPos.Y + sz
	return image.Rect(p.X, p.Y, maxX, maxY)
}

func (b *Building) Update() error {
	return nil
}

func (b *Building) Render(renderer platform.Renderer, cameraPos image.Point) error {
	img := b.buildingSprite.Imgs[0]
	renderer.BlitImage(b.worldPos.Add(cameraPos), img, b.palette)
	return nil
}
