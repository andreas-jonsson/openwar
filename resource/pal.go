/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
There are three palette formats used by different versions of the game. The DOS versions
palettes have no header and are to be read as RGB byte values between 0 and 64, each,
so they need to be multiplied by 4 to get full 8 bit values for each channel. Those DOS
palettes can also be halved and contain only 128 colors instead of 256.

Mac palettes always have 256 colors. They begin with an 8 byte header, followed by
a list of all colors. Each entry in that list consists of the color index as a 2 byte integer,
followed by RRGGBB, with each color value written twice, so one byte for each channel
can be ignored. An image will typically require two palettes, one for the lower half of
indices and one for the greater half. The lower half is usually a terrain palette.
*/

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
