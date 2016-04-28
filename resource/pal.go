/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package resource

import (
	"encoding/binary"
	"errors"
	"image/color"
	"path"
)

type Palettes map[string]color.Palette

func LoadPal(arch *Archive) (Palettes, error) {
	palettes := make(map[string]color.Palette)

	for file, data := range arch.Files {
		if path.Ext(file) == ".PAL" {
			paletteSize := len(data) / 3

			if paletteSize != 128 && paletteSize != 256 {
				return palettes, errors.New("invalid palette")
			}

			pal := make([]color.Color, paletteSize)
			reader, _ := arch.Open(file)

			var rgb [3]byte
			for i := 0; i < paletteSize; i++ {
				if err := binary.Read(reader, binary.LittleEndian, &rgb); err != nil {
					return palettes, err
				}
				pal[i] = color.RGBA{rgb[0] * 4, rgb[1] * 4, rgb[2] * 4, 255}
			}

			palettes[file] = pal
		}
	}

	return palettes, nil
}
