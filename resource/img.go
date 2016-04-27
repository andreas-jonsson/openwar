/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package resource

import (
	"encoding/binary"
	"image"
	"image/color/palette"
	"path"
)

type Images map[string]*image.Paletted

func LoadImg(arch *Archive) (Images, error) {
	images := make(map[string]*image.Paletted)

	for file, data := range arch.Files {
		if path.Ext(file) == ".IMG" {
			var (
				width, height uint16
			)

			reader, _ := arch.Open(file)
			if err := binary.Read(reader, binary.LittleEndian, &width); err != nil {
				return images, err
			}

			if err := binary.Read(reader, binary.LittleEndian, &height); err != nil {
				return images, err
			}

			img := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), palette.Plan9)
			copy(img.Pix, data[4:])
			images[file] = img
		}
	}

	return images, nil
}
