/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package resource

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"path"
	"strings"
)

type (
	Tileset struct {
		NumTiles int
		Data     *image.Paletted
		Palette  color.Palette
	}

	Tilesets map[string]Tileset
)

func LoadTilesets(arch *Archive, images Images, palettes Palettes) (Tilesets, error) {
	tilesets := make(Tilesets)

	for file, data := range arch.Files {
		ext := path.Ext(file)
		if ext == ".TIL" {
			megaTileData, ok := arch.Files[strings.TrimSuffix(file, ext)+".PTR"]
			if !ok {
				fmt.Fprintf(Logger, "%s is incomplete, missing megatile.\n", file)
				continue
			}

			pal, ok := palettes[strings.TrimSuffix(file, ext)+".PAL"]
			if !ok {
				fmt.Fprintf(Logger, "%s is incomplete, missing palette.\n", file)
				continue
			}

			reader := bytes.NewReader(megaTileData)
			numTiles := len(megaTileData) / 8

			megaTile := make([]uint16, len(megaTileData)/2)
			if err := binary.Read(reader, binary.LittleEndian, &megaTile); err != nil {
				return tilesets, err
			}

			img := image.NewPaletted(image.Rect(0, 0, numTiles*16, 16), pal)
			for i := 0; i < numTiles; i++ {
				makeTile(img, i, data, megaTile[i*4:])
			}

			images[file] = Image{Data: img}
			tilesets[file] = Tileset{Data: img, Palette: pal, NumTiles: numTiles}
		}
	}

	return tilesets, nil
}

func makeTile(img *image.Paletted, tileIndex int, miniTile []byte, megaTile []uint16) {
	var (
		offset, index int
		flipX, flipY  bool
		tile          [256]byte

		flip = []int{7, 6, 5, 4, 3, 2, 1, 0, 8}
	)

	for t := 0; t < 2; t++ {
		for s := 0; s < 2; s++ {
			offset = int(megaTile[s+t*2])
			flipX = offset&2 != 0
			flipY = offset&1 != 0
			index = (offset & 0xFFFC) << 1

			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					fy := y
					fx := x

					if flipY {
						fy = flip[y]
					}
					if flipX {
						fx = flip[x]
					}

					tile[128*t+16*y+8*s+x] = miniTile[index+fy*8+fx]
				}
			}
		}
	}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.SetColorIndex(tileIndex*16+x, y, tile[y*16+x])
		}
	}
}
