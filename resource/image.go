/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
Image files start with a 2 byte width and a 2 byte height, followed by the
pixels of the image as 1 byte color indices from the corresponding palette.
*/

package resource

import (
	"encoding/binary"
	"image"
	"image/color/palette"
	"io"
	"path"
)

type Images map[string]*image.Paletted

func LoadImages(arch *Archive) (Images, error) {
	images := make(map[string]*image.Paletted)

	for file, data := range arch.Files {
		if path.Ext(file) == ".IMG" {
			reader, _ := arch.Open(file)

			img, err := loadImageData(reader, data[4:])
			if err != nil {
				return images, err
			}

			images[file] = img
		}
	}

	return images, nil
}

func loadImageData(reader io.Reader, pix []byte) (*image.Paletted, error) {
	var (
		img           *image.Paletted
		width, height uint16
	)

	if err := binary.Read(reader, binary.LittleEndian, &width); err != nil {
		return img, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &height); err != nil {
		return img, err
	}

	img = image.NewPaletted(image.Rect(0, 0, int(width), int(height)), palette.Plan9)
	copy(img.Pix, pix)
	return img, nil
}
