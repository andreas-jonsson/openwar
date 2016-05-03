/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"encoding/binary"
	"image"
	"image/color"
	"math"

	"github.com/openwar-hq/openwar/resource"
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

func newTerrain(g *Game) (*terrain, error) {
	ter := &terrain{g: g}
	ter.tileset = g.resources.Tilesets["SWAMP.TIL"]
	ter.pal = g.resources.Palettes["SWAMP.PAL"]

	file := "ORC01.TER"
	reader, err := g.resources.Archive.Open(file)
	if err != nil {
		return nil, err
	}

	size := len(g.resources.Archive.Files[file]) / 2
	ter.tileIndex = make([]uint16, size)
	ter.mapSize = int(math.Sqrt(float64(size)))

	if err := binary.Read(reader, binary.LittleEndian, ter.tileIndex); err != nil {
		return nil, err
	}

	return ter, nil
}

func (ter *terrain) render(cullRect image.Rectangle) {
	renderer := ter.g.renderer
	min := cullRect.Min
	max := cullRect.Max

	for y := min.Y; y < ter.mapSize && y < max.Y; y++ {
		for x := min.X; x < ter.mapSize && x < max.X; x++ {
			idx := int(ter.tileIndex[y*ter.mapSize+x])
			if idx > ter.tileset.NumTiles-1 {
				panic("index out of range")
			}

			rect := image.Rect(16*idx, 0, 16*idx+16, 16)
			src := ter.tileset.Data
			tilePos := image.Point{x * 16, y * 16}

			renderer.Blit(tilePos, src, rect, ter.pal)
		}
	}
}
