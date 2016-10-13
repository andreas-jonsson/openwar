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
	"image/draw"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/andreas-jonsson/openwar/resource"
)

type environmentType int

const (
	environmentForest environmentType = iota
	environmentSwamp
	environmentDungeon
)

const paletteAnimationSpeed = 250 * time.Millisecond

type (
	terrainImpl struct {
		g *Game

		tileset resource.Tileset
		pal     color.Palette
		palAnim *time.Ticker

		mapSize   int
		tileIndex []uint16
		tileFlags []uint16

		miniMap  *image.RGBA
		mapImage *image.Paletted
	}

	terrain interface {
		render(cullRect image.Rectangle, cameraPos image.Point)
		miniMapImage() *image.RGBA
		terrainPalette() color.Palette
		size() int
	}
)

// HUMAN04, HUMAN08 and CUSTOMD3-8 is not working...

var mapsEnvironment = map[string]environmentType{
	"CUSTOMD1": environmentDungeon,
	"CUSTOMD2": environmentDungeon,
	"CUSTOMD3": environmentDungeon,
	"CUSTOMD4": environmentDungeon,
	"CUSTOMD5": environmentDungeon,
	"CUSTOMD6": environmentDungeon,
	"CUSTOMD7": environmentDungeon,
	"CUSTOMD8": environmentDungeon,
	"CUSTOMF1": environmentForest,
	"CUSTOMF2": environmentForest,
	"CUSTOMS1": environmentSwamp,
	"CUSTOMS2": environmentSwamp,
	"HUMAN01":  environmentForest,
	"HUMAN02":  environmentForest,
	"HUMAN03":  environmentSwamp,
	"HUMAN04":  environmentForest,
	"HUMAN05":  environmentForest,
	"HUMAN06":  environmentForest,
	"HUMAN07":  environmentForest,
	"HUMAN08":  environmentSwamp,
	"HUMAN09":  environmentSwamp,
	"HUMAN10":  environmentSwamp,
	"HUMAN11":  environmentSwamp,
	"HUMAN12":  environmentSwamp,
	"ORC01":    environmentSwamp,
	"ORC02":    environmentSwamp,
	"ORC03":    environmentSwamp,
	"ORC04":    environmentSwamp,
	"ORC05":    environmentSwamp,
	"ORC06":    environmentSwamp,
	"ORC07":    environmentSwamp,
	"ORC08":    environmentSwamp,
	"ORC09":    environmentSwamp,
	"ORC10":    environmentSwamp,
	"ORC11":    environmentSwamp,
	"ORC12":    environmentSwamp,
}

func newTerrain(g *Game, name string) (terrain, error) {
	ter := &terrainImpl{g: g}

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
	ter.palAnim = time.NewTicker(paletteAnimationSpeed)

	reader, err := g.resources.Archive.Open(name + ".TER")
	if err != nil {
		return nil, err
	}

	size := len(g.resources.Archive.Files[name+".TER"]) / 2
	ter.tileIndex = make([]uint16, size)
	ter.mapSize = int(math.Sqrt(float64(size)))

	if err = binary.Read(reader, binary.LittleEndian, ter.tileIndex); err != nil {
		return nil, err
	}

	reader, err = g.resources.Archive.Open(name + ".FOG")
	if err != nil {
		return nil, err
	}

	ter.tileFlags = make([]uint16, size)
	if err := binary.Read(reader, binary.LittleEndian, ter.tileFlags); err != nil {
		return nil, err
	}

	ter.mapImage, ter.miniMap = ter.createMap()
	return ter, nil
}

func (ter *terrainImpl) terrainPalette() color.Palette {
	return ter.pal
}

func (ter *terrainImpl) miniMapImage() *image.RGBA {
	return ter.miniMap
}

func (ter *terrainImpl) size() int {
	return ter.mapSize
}

func (ter *terrainImpl) render(cullRect image.Rectangle, cameraPos image.Point) {
	cullMin := cullRect.Min
	cullRect.Max = cullRect.Size()
	cullRect.Min = image.ZP
	cullRect = cullRect.Add(cameraPos)

	ter.animatePalette()
	ter.g.renderer.Blit(cullMin, ter.mapImage, cullRect, ter.pal)
}

func (ter *terrainImpl) animatePalette() {
	const (
		startFrame = 112
		numFrames  = 7
	)

	select {
	case <-ter.palAnim.C:
		pal := ter.pal
		tail := pal[startFrame+numFrames]
		copy(pal[startFrame+1:startFrame+numFrames+1], pal[startFrame:startFrame+numFrames])
		pal[startFrame] = tail
	default:
	}
}

func (ter *terrainImpl) createMap() (*image.Paletted, *image.RGBA) {
	miniMap := image.NewRGBA(image.Rect(0, 0, ter.mapSize, ter.mapSize))
	mapImage := image.NewPaletted(image.Rect(0, 0, ter.mapSize*16, ter.mapSize*16), ter.pal)

	for y := 0; y < ter.mapSize; y++ {
		for x := 0; x < ter.mapSize; x++ {
			offset := y*ter.mapSize + x
			idx := int(ter.tileIndex[offset])

			if idx > ter.tileset.NumTiles-1 {
				log.Panicln("index out of range", idx, ter.tileset.NumTiles-1)
			}

			// Render minimap.
			miniMap.Set(x, y, tileColor(ter.tileset.Data, ter.pal, idx))

			rect := image.Rect(0, 16*idx, 16, 16*idx+16)
			src := ter.tileset.Data
			tilePos := image.Point{x * 16, y * 16}

			r := image.Rectangle{tilePos, tilePos.Add(rect.Size())}
			draw.Draw(mapImage, r, src, rect.Min, draw.Src)

			flags := ter.tileFlags[offset]
			if flags != 0 {
				rect = image.Rect(tilePos.X, tilePos.Y, tilePos.X+16, tilePos.Y+16)
				ter.g.renderer.DrawRect(rect, color.RGBA{byte(flags & 0xFF), 0x0, 0x0, 0xFF}, false)
			}
		}
	}
	return mapImage, miniMap
}

func tileColor(img *image.Paletted, pal color.Palette, idx int) color.Color {
	return pal[img.ColorIndexAt(rand.Intn(7)+4, idx*16+rand.Intn(7)+4)]
}
