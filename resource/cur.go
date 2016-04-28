/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
Cursors start with x and y offset, each as 2 byte integer,
followed by the usual structure of an .IMG file.
*/

package resource

import (
	"encoding/binary"
	"image"
	"path"
)

type Cursor struct {
	X, Y int
	Img  *image.Paletted
}

type Cursors map[string]Cursor

func LoadCur(arch *Archive, images Images) (Cursors, error) {
	cursors := make(map[string]Cursor)

	for file, data := range arch.Files {
		if path.Ext(file) == ".CUR" {
			var (
				x, y uint16
			)

			reader, _ := arch.Open(file)
			if err := binary.Read(reader, binary.LittleEndian, &x); err != nil {
				return cursors, err
			}

			if err := binary.Read(reader, binary.LittleEndian, &y); err != nil {
				return cursors, err
			}

			img, err := loadImgData(reader, data[8:])
			if err != nil {
				return cursors, err
			}

			images[file] = img
			cursors[file] = Cursor{X: int(x), Y: int(y), Img: img}
		}
	}

	return cursors, nil
}
