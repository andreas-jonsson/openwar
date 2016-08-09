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
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/andreas-jonsson/openwar/resource"
)

type environmentType int

const (
	environmentForest environmentType = iota
	environmentSwamp
	environmentDungeon
)

type terrain struct {
	g *Game

	mapSize   int
	tileIndex []uint16
	tileset   resource.Tileset
	pal       color.Palette
}

var mapsEnvironment = map[string]environmentType{
	"ORC01.TER": environmentSwamp,

	"HUMAN01.TER": environmentForest,
}

func newTerrain(g *Game, name string) (*terrain, error) {
	ter := &terrain{g: g}

	env, ok := mapsEnvironment[name]
	if !ok {
		return nil, errors.New("invalid terrain name")
	}

	switch env {
	case environmentForest:
		ter.tileset = g.resources.Tilesets["FOREST.TIL"]
		ter.pal = g.resources.Palettes["FOREST.PAL"]
	case environmentSwamp:
		ter.tileset = g.resources.Tilesets["SWAMP.TIL"]
		ter.pal = g.resources.Palettes["SWAMP.PAL"]
	case environmentDungeon:
		ter.tileset = g.resources.Tilesets["DUNGEON.TIL"]
		ter.pal = g.resources.Palettes["DUNGEON.PAL"]
	}

	reader, err := g.resources.Archive.Open(name)
	if err != nil {
		return nil, err
	}

	size := len(g.resources.Archive.Files[name]) / 2
	ter.tileIndex = make([]uint16, size)
	ter.mapSize = int(math.Sqrt(float64(size)))

	if err := binary.Read(reader, binary.LittleEndian, ter.tileIndex); err != nil {
		return nil, err
	}

	return ter, nil
}

func (ter *terrain) render(cullRect image.Rectangle, cameraPos image.Point) {
	renderer := ter.g.renderer

	min := cullRect.Min.Add(cameraPos).Div(16)
	max := cullRect.Max.Add(cameraPos).Div(16)

	max.X++
	max.Y++

	for y, dy := min.Y, 0; y < ter.mapSize && y < max.Y; y++ {
		for x, dx := min.X, 0; x < ter.mapSize && x < max.X; x++ {
			idx := int(ter.tileIndex[y*ter.mapSize+x])
			if idx > ter.tileset.NumTiles-1 {
				panic("index out of range")
			}

			rect := image.Rect(0, 16*idx, 16, 16*idx+16)
			src := ter.tileset.Data
			tilePos := image.Point{dx*16 - (cameraPos.X % 16), dy*16 - (cameraPos.Y % 16)}

			renderer.Blit(tilePos, src, rect, ter.pal)
			dx++
		}
		dy++
	}
}
