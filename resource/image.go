/*
Copyright (C) 2016-2018 Andreas T Jonsson

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

/*
Image files start with a 2 byte width and a 2 byte height, followed by the
pixels of the image as 1 byte color indices from the corresponding palette.

Cursors start with x and y offset, each as 2 byte integer, followed by the
usual structure of an .IMG file.
*/

package resource

import (
	"encoding/binary"
	"image"
	"image/color/palette"
	"io"
	"path"
)

type (
	Image struct {
		X, Y int
		Data *image.Paletted
	}

	Images map[string]Image
)

func (img *Image) Point() image.Point {
	return image.Point{img.X, img.Y}
}

func LoadImages(arch *Archive) (Images, error) {
	images := make(Images)

	for file, data := range arch.Files {
		var (
			err       error
			x, y      uint16
			imageData *image.Paletted
		)

		switch path.Ext(file) {
		case ".IMG":
			reader, _ := arch.Open(file)

			imageData, err = loadImageData(reader, data[4:])
			if err != nil {
				return images, err
			}
		case ".CUR":
			reader, _ := arch.Open(file)

			if err = binary.Read(reader, binary.LittleEndian, &x); err != nil {
				return images, err
			}

			if err = binary.Read(reader, binary.LittleEndian, &y); err != nil {
				return images, err
			}

			imageData, err = loadImageData(reader, data[8:])
			if err != nil {
				return images, err
			}
		}

		if imageData != nil {
			images[file] = Image{X: int(x), Y: int(y), Data: imageData}
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
