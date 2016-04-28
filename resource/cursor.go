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

type (
	Cursor struct {
		X, Y uint16
		Img  *image.Paletted
	}

	Cursors map[string]Cursor
)

func LoadCursors(arch *Archive, images Images) (Cursors, error) {
	cursors := make(map[string]Cursor)

	for file, data := range arch.Files {
		if path.Ext(file) == ".CUR" {
			reader, _ := arch.Open(file)

			var cursor Cursor
			if err := binary.Read(reader, binary.LittleEndian, &cursor.X); err != nil {
				return cursors, err
			}

			if err := binary.Read(reader, binary.LittleEndian, &cursor.Y); err != nil {
				return cursors, err
			}

			img, err := loadImageData(reader, data[8:])
			if err != nil {
				return cursors, err
			}

			cursor.Img = img
			images[file] = img
			cursors[file] = cursor
		}
	}

	return cursors, nil
}
